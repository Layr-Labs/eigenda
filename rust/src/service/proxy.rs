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

#[cfg(test)]
pub mod tests {
    use std::{borrow::Cow, collections::HashMap};

    use testcontainers::{
        ContainerAsync, Image,
        core::{ContainerPort, WaitFor},
        runners::AsyncRunner,
    };

    use crate::service::proxy::ProxyClient;

    pub async fn create_test_eigenda_proxy()
    -> Result<(ProxyClient, ContainerAsync<EigenDaProxy>), anyhow::Error> {
        let container = EigenDaProxy::default().start().await?;
        let host_port = container.get_host_port_ipv4(PORT).await?;
        let proxy = ProxyClient::new(format!("http://127.0.0.1:{host_port}"))?;

        Ok((proxy, container))
    }

    #[tokio::test]
    async fn blob_roundtrip() {
        let (proxy, _container) = create_test_eigenda_proxy().await.unwrap();
        let blob = vec![0; 1000];

        // Store the blob
        let certificate = proxy.store_blob(&blob).await.unwrap();

        // Retrieve the blob
        let retrieved_blob = proxy.get_blob(&certificate).await.unwrap();

        assert_eq!(blob, retrieved_blob);
    }

    const NAME: &str = "ghcr.io/layr-labs/eigenda-proxy";
    const TAG: &str = "latest";
    const READY_MSG: &str = "Started EigenDA proxy server";
    const PORT: ContainerPort = ContainerPort::Tcp(3100);

    /// EigenDAProxy image for testcontainers
    #[derive(Debug)]
    pub struct EigenDaProxy {
        env_vars: HashMap<String, String>,
    }

    impl Default for EigenDaProxy {
        fn default() -> Self {
            let mut env_vars = HashMap::new();
            env_vars.insert("EIGENDA_PROXY_PORT".to_owned(), PORT.as_u16().to_string());
            env_vars.insert(
                "EIGENDA_PROXY_MEMSTORE_ENABLED".to_owned(),
                "true".to_string(),
            );

            Self { env_vars }
        }
    }

    impl Image for EigenDaProxy {
        fn name(&self) -> &str {
            NAME
        }

        fn tag(&self) -> &str {
            TAG
        }

        fn ready_conditions(&self) -> Vec<WaitFor> {
            vec![WaitFor::message_on_stdout(READY_MSG)]
        }

        fn env_vars(
            &self,
        ) -> impl IntoIterator<Item = (impl Into<Cow<'_, str>>, impl Into<Cow<'_, str>>)> {
            &self.env_vars
        }

        fn expose_ports(&self) -> &[ContainerPort] {
            &[PORT]
        }
    }
}
