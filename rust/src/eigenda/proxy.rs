use std::{str::FromStr, time::Duration};

use backon::{ExponentialBuilder, Retryable};
use bytes::Bytes;
use hex::encode;
use reqwest::{Request, Url, header::CONTENT_TYPE};
use thiserror::Error;
use tracing::{error, trace};

use crate::{
    eigenda::cert::{StandardCommitment, StandardCommitmentParseError},
    service::config::EigenDaConfig,
};

/// Default maximal number of times we retry requests.
const DEFAULT_MAX_RETRY_TIMES: u64 = 10;
/// Default starting delay at which requests will be retried.
const DEFAULT_MIN_RETRY_DELAY: Duration = Duration::from_millis(1000);
/// Default maximal delay at which requests will be retried.
const DEFAULT_MAX_RETRY_DELAY: Duration = Duration::from_secs(10);

#[derive(Debug, Clone)]
pub struct ProxyClient {
    url: Url,
    inner: reqwest::Client,
    // Backoff for retrying strategy
    backoff: Option<ExponentialBuilder>,
}

impl ProxyClient {
    pub fn new(config: &EigenDaConfig) -> Result<Self, ProxyError> {
        let min_retry_delay = config
            .proxy_min_retry_delay
            .map(Duration::from_millis)
            .unwrap_or(DEFAULT_MIN_RETRY_DELAY);

        let max_retry_delay = config
            .proxy_max_retry_delay
            .map(Duration::from_millis)
            .unwrap_or(DEFAULT_MAX_RETRY_DELAY);

        let max_retry_times = config
            .proxy_max_retry_times
            .unwrap_or(DEFAULT_MAX_RETRY_TIMES);

        let backoff = ExponentialBuilder::default()
            .with_min_delay(min_retry_delay)
            .with_max_delay(max_retry_delay)
            .with_max_times(max_retry_times as usize);

        let url = Url::from_str(&config.proxy_url)?;
        let inner = reqwest::Client::builder().build()?;

        Ok(Self {
            url,
            inner,
            backoff: Some(backoff),
        })
    }

    /// Fetch blob data for the given certificate.
    pub async fn get_blob(&self, certificate: &StandardCommitment) -> Result<Bytes, ProxyError> {
        let hex = encode(certificate.to_rlp_bytes());
        let mut url = self.url.join(&format!("/get/0x{hex}"))?;
        url.set_query(Some("commitment_mode=standard"));

        let request = self.inner.get(url).build()?;
        let response = self.call(request).await?;
        Ok(response)
    }

    /// Stores the new blob and returns a certificate
    pub async fn store_blob(&self, blob: &[u8]) -> Result<StandardCommitment, ProxyError> {
        let mut url = self.url.join("/put")?;
        url.set_query(Some("commitment_mode=standard"));

        let request = self
            .inner
            .post(url)
            .header(CONTENT_TYPE, "application/octet-stream")
            .body(blob.to_vec())
            .build()?;

        let response = self.call(request).await?;

        // We optimistically expect a certificate
        match StandardCommitment::from_rlp_bytes(response.as_ref()) {
            Ok(cert) => Ok(cert),
            Err(err) => {
                // Try to serialize a string from response bytes for a nicer
                // error message. This error handling could be better. But
                // currently, we don't really know what to expect from the proxy
                // in cases when response is not a cert
                match str::from_utf8(&response) {
                    Ok(response_body) => {
                        error!(?err, %response_body, "Error occurred while parsing proxy response");
                    }
                    Err(_) => {
                        error!(
                            ?err,
                            response_body = ?response,
                            "Error occurred while parsing proxy response"
                        );
                    }
                }

                Err(err.into())
            }
        }
    }

    async fn call(&self, request: Request) -> Result<Bytes, reqwest::Error> {
        // If there is retry strategy, run with retries, otherwise just call once
        if let Some(backoff) = self.backoff.as_ref() {
            // The operation to be retried
            let request = &request;
            let operation = || async {
                let request = request
                    .try_clone()
                    .expect("the body is not a stream. so the request is clone-able");
                self.call_inner(request).await
            };

            // Notification on each retry
            let notify = |err: &reqwest::Error, dur: Duration| trace!(?request, ?dur, %err, "eigenda proxy error");

            operation
                .retry(backoff)
                .when(|err| err.is_connect() || err.is_timeout())
                .notify(notify)
                .await
        } else {
            self.call_inner(request).await
        }
    }

    async fn call_inner(&self, request: Request) -> Result<Bytes, reqwest::Error> {
        let request = self.inner.execute(request).await?;
        let bytes = request.bytes().await?;

        Ok(bytes)
    }
}

/// Represents errors that can occur during EigenDA proxy operations.
#[derive(Debug, Error)]
pub enum ProxyError {
    /// Error when parsing URL.
    #[error("Url parse error: {0}")]
    UrlParse(#[from] url::ParseError),

    /// Error when sending an HTTP request.
    #[error("HTTP error: {0}")]
    Http(reqwest::Error),

    /// Error when the HTTP request times out.
    #[error("HTTP request timed out")]
    HttpTimeout,

    /// Error parsing the commitment
    #[error("StandardCommitmentParseError: {0}")]
    StandardCommitmentParseError(#[from] StandardCommitmentParseError),
}

impl From<reqwest::Error> for ProxyError {
    fn from(error: reqwest::Error) -> Self {
        match error {
            error if error.is_timeout() => ProxyError::HttpTimeout,
            error => ProxyError::Http(error),
        }
    }
}
