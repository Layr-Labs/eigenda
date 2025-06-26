use hex::encode;
use reqwest::{IntoUrl, Url, header::CONTENT_TYPE};
use thiserror::Error;

#[derive(Debug, Clone)]
pub struct ProxyClient {
    inner: reqwest::Client,
    base_url: Url,
}

impl ProxyClient {
    pub fn new<U>(url: U) -> Result<Self, ProxyError>
    where
        U: IntoUrl,
    {
        let inner = reqwest::Client::builder().build()?;

        Ok(Self {
            inner,
            base_url: url.into_url()?,
        })
    }

    /// Fetch blob data for the given certificate
    pub async fn get_blob(&self, cert: &[u8]) -> Result<Vec<u8>, ProxyError> {
        let hex = encode(cert);
        let mut url = self.base_url.join(&format!("/get/0x{hex}"))?;
        url.set_query(Some("commitment_mode=standard"));

        let request = self.inner.get(url).build()?;

        let response = self.inner.execute(request).await?;
        let response = response.bytes().await?;

        Ok(response.to_vec())
    }

    /// Stores the new blob and returns a certificate
    pub async fn store_blob(&self, blob: &[u8]) -> Result<Vec<u8>, ProxyError> {
        let mut url = self.base_url.join("/put")?;
        url.set_query(Some("commitment_mode=standard"));

        let request = self
            .inner
            .post(url)
            .header(CONTENT_TYPE, "application/octet-stream")
            .body(blob.to_vec())
            .build()?;

        let response = self.inner.execute(request).await?;
        let response = response.bytes().await?;

        Ok(response.to_vec())
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
}

impl From<reqwest::Error> for ProxyError {
    fn from(error: reqwest::Error) -> Self {
        match error {
            error if error.is_timeout() => ProxyError::HttpTimeout,
            error => ProxyError::Http(error),
        }
    }
}
