use tracing::error;
use tracing_subscriber::{EnvFilter, Registry, layer::SubscriberExt};
use tracing_tree::HierarchicalLayer;

pub fn init_tracing() {
    let subscriber = Registry::default()
        .with(EnvFilter::from_default_env())
        .with(HierarchicalLayer::new(2).with_indent_amount(4));
    if let Err(err) = tracing::subscriber::set_global_default(subscriber) {
        error!("{err}");
    }
}
