use std::str::FromStr;
use std::time::Duration;

use backon::{ExponentialBuilder, Retryable};
use bytes::Bytes;
use hex::encode;
use reqwest::header::CONTENT_TYPE;
use reqwest::{Request, Url};
use thiserror::Error;
use tracing::{error, trace};

use crate::eigenda::cert::{StandardCommitment, StandardCommitmentParseError};
use crate::service::config::EigenDaConfig;

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

    /// Fetch encoded payload data for the given certificate.
    pub async fn get_encoded_payload(
        &self,
        certificate: &StandardCommitment,
    ) -> Result<Bytes, ProxyError> {
        let hex = encode(certificate.to_rlp_bytes());
        let mut url = self.url.join(&format!("/get/0x{hex}"))?;
        url.set_query(Some("commitment_mode=standard&return_encoded_payload=true"));

        let request = self.inner.get(url).build()?;
        let response = self.call(request).await?;
        Ok(response)
    }

    /// Stores the payload and returns a certificate
    pub async fn store_payload(&self, payload: &[u8]) -> Result<StandardCommitment, ProxyError> {
        let mut url = self.url.join("/put")?;
        url.set_query(Some("commitment_mode=standard"));

        let request = self
            .inner
            .post(url)
            .header(CONTENT_TYPE, "application/octet-stream")
            .body(payload.to_vec())
            .build()?;

        let response = self.call(request).await?;

        // We optimistically expect a certificate
        match StandardCommitment::from_rlp_bytes(response.as_ref()) {
            Ok(cert) => Ok(cert),
            Err(err) => {
                let response = str::from_utf8(&response);
                error!(
                    ?err,
                    ?response,
                    "Error occurred while parsing proxy response"
                );

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

#[cfg(test)]
mod tests {
    use wiremock::matchers::{header, method, path, query_param};
    use wiremock::{Mock, MockServer, ResponseTemplate};

    use super::*;
    use crate::service::config::{EigenDaConfig, Network};

    fn create_test_config(proxy_url: String) -> EigenDaConfig {
        EigenDaConfig {
            network: Network::Holesky,
            ethereum_rpc_url: "http://test.com".to_string(),
            ethereum_compute_units: None,
            ethereum_max_retry_times: None,
            ethereum_initial_backoff: None,
            proxy_url,
            proxy_min_retry_delay: Some(100),
            proxy_max_retry_delay: Some(1000),
            proxy_max_retry_times: Some(3),
            sequencer_signer: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
                .to_string(),
        }
    }

    fn create_test_certificate() -> StandardCommitment {
        let commitment_hex = "02f90389e5a0c769488dd5264b3ef21dce7ee2d42fba43e1f83ff228f501223e38818cb14492833f44fcf901eff901caf9018180820001f90159f842a0012e810ffc0a83074b3d14db9e78bbae623f7770cac248df9e73fac6b9d59d17a02a916ffbbf9dde4b7ebe94191a29ff686422d7dcb3b47ecb03c6ada75a9c15c8f888f842a01811c8b4152fce9b8c4bae61a3d097e61dfc43dc7d45363d19e7c7f1374034ffa001edc62174217cdce60a4b52fa234ac0d96db4307dac9150e152ba82cbb4d2f1f842a00f423b0dbc1fe95d2e3f7dbac6c099e51dbf73400a4b3f26b9a29665b4ac58a8a01855a2bd56c0e8f4cc85ac149cf9a531673d0e89e22f0d6c4ae419ed7c5d2940f888f842a02667cbb99d60fa0d7f3544141d3d531dceeeb50b06e5a0cdc42338a359138ae4a00dff4c929d8f8a307c19bba6e8006fe6700f6554cef9eb3797944f89472ffb30f842a004c17a6225acd5b4e7d672a1eb298c5358f4f6f17d04fd1ee295d0c0d372fa84a024bc3ad4d5e54f54f71db382ce276f37ac3c260cc74306b832e8a3c93c7951d302a0e43e11e2405c2fd1d880af8612d969b654827e0ba23d9feb3722ccce6226fce7b8411ddf4553c79c0515516fd3c8b3ae6a756b05723f4d0ebe98a450c8bcc96cbb355ef07a44eeb56f831be73647e4da20e22fa859f984ee41d6efcd3692063b0b0601c2800101a0a69e552a6fc2ff75d32edaf5313642ddeebe60d2069435d12e266ce800e9e96bf9016bc0c0f888f842a00d45727a99053af8d38d4716ab83ace676096e7506b6b7aa6953e87bc04a023ca016c030c31dd1c94062948ecdce2e67c4e6626c16af0033dcdb7a96362c937d48f842a00a95fac74aba7e3fbd24bc62457ce6981803d8f5fef28871d3d5e2af05d50cd4a0117400693917cd50d9bc28d4ab4fadf93a23e771f303637f8d1f83cd0632c3fcf888f842a0301bfced3253e99e8d50f2fed62313a16d714013d022a4dc4294656276f10d1ba0152e047a83c326a9d81dac502ec429b662b58ee119ca4c8748a355b539c24131f842a01944b5b4a3e93d46b0fe4370128c6cdcd066ae6b036b019a20f8d22fe9a10d67a00ddf3421722967c0bd965b9fc9e004bf01183b6206fec8de65e40331d185372ef842a02db8fb278708abf8878ebf578872ab35ee914ad8196b78de16b34498222ac1c2a02ff9d9a5184684f4e14530bde3a61a2f9adaa74734dff104b61ba3d963a644dac68207388208b7c68209998209c5c2c0c0820001";
        let raw_commitment = hex::decode(commitment_hex).expect("Valid test certificate hex");
        StandardCommitment::from_rlp_bytes(raw_commitment.as_slice())
            .expect("Valid test certificate")
    }

    #[tokio::test]
    async fn test_get_encoded_payload_success() {
        let mock_server = MockServer::start().await;
        let config = create_test_config(mock_server.uri());
        let client = ProxyClient::new(&config).unwrap();

        let test_data = b"test encoded payload data";
        let certificate = create_test_certificate();
        let hex_cert = hex::encode(certificate.to_rlp_bytes());

        Mock::given(method("GET"))
            .and(path(format!("/get/0x{hex_cert}")))
            .and(query_param("commitment_mode", "standard"))
            .and(query_param("return_encoded_payload", "true"))
            .respond_with(ResponseTemplate::new(200).set_body_bytes(test_data))
            .mount(&mock_server)
            .await;

        let payload = client.get_encoded_payload(&certificate).await.unwrap();
        assert_eq!(payload.as_ref(), test_data);
    }

    #[tokio::test]
    async fn test_get_encoded_payload_http_error() {
        let mock_server = MockServer::start().await;
        let mut config = create_test_config(mock_server.uri());
        // Disable retries for this test to ensure error propagation
        config.proxy_max_retry_times = Some(0);
        let mut client = ProxyClient::new(&config).unwrap();
        client.backoff = None;

        let certificate = create_test_certificate();
        let hex_cert = hex::encode(certificate.to_rlp_bytes());

        Mock::given(method("GET"))
            .and(path(format!("/get/0x{hex_cert}")))
            .respond_with(ResponseTemplate::new(500).set_body_string("Internal Server Error"))
            .mount(&mock_server)
            .await;

        let payload = client.get_encoded_payload(&certificate).await.unwrap();
        assert_eq!(payload.as_ref(), b"Internal Server Error");
    }

    #[tokio::test]
    async fn test_store_payload_success() {
        let mock_server = MockServer::start().await;
        let config = create_test_config(mock_server.uri());
        let client = ProxyClient::new(&config).unwrap();

        let test_payload = b"test payload to store";
        let certificate = create_test_certificate();
        let cert_rlp_bytes = certificate.to_rlp_bytes();

        Mock::given(method("POST"))
            .and(path("/put"))
            .and(query_param("commitment_mode", "standard"))
            .and(header("content-type", "application/octet-stream"))
            .respond_with(ResponseTemplate::new(200).set_body_bytes(cert_rlp_bytes.as_ref()))
            .mount(&mock_server)
            .await;

        let returned_cert = client.store_payload(test_payload).await.unwrap();
        assert_eq!(returned_cert.to_rlp_bytes(), cert_rlp_bytes);
    }

    #[tokio::test]
    async fn test_store_payload_invalid_certificate_response() {
        let mock_server = MockServer::start().await;
        let config = create_test_config(mock_server.uri());
        let client = ProxyClient::new(&config).unwrap();

        let test_payload = b"test payload to store";

        Mock::given(method("POST"))
            .and(path("/put"))
            .and(query_param("commitment_mode", "standard"))
            .respond_with(ResponseTemplate::new(200).set_body_string("invalid certificate data"))
            .mount(&mock_server)
            .await;

        let err = client.store_payload(test_payload).await.unwrap_err();
        assert!(matches!(err, ProxyError::StandardCommitmentParseError(_)));
    }
}
