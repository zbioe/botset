//! Requires chromedriver running on port 4444:
//!
//!     chromedriver --port=4444
//!
//! Run as follows:
//!
//!     cargo run --example minimal_async

use std::time::Duration;
use thirtyfour::prelude::*;
use thirtyfour::query::*;

#[tokio::main]
async fn main() -> WebDriverResult<()> {
    let mut caps = DesiredCapabilities::chrome();
    caps.add_chrome_option(
        "prefs",
        serde_json::json!({
            "intl.accept_languages": "en,en_US"
        }),
    )?;

    let driver = WebDriver::new("http://localhost:4444", caps).await?;
    let poller =
        ElementPoller::TimeoutWithInterval(Duration::new(10, 0), Duration::from_millis(500));
    driver.set_query_poller(poller);
    // Navigate to twitter
    driver.get("https://twitter.com/i/flow/signup").await?;
    // sleep
    // let elem = ElementWaitable::wait_until(&driver);
    // elem.wait_until().exists().await?;
    // let _elem = driver
    //     .query(By::XPath("//div[@role='button' and contains(.,'mail')]"))
    //     .wait_until()
    //     .exists()
    // .await?;
    //elem.wait_until().displayed().await?;
    //let _is_found = driver.query(By::Id("button1")).nowait().exists().await?;

    // click
    let elem_button = driver
        .query(By::XPath("//div[@role='button' and contains(.,'mail')]"))
        .single()
        .await?;
    elem_button.click().await?;

    let elem_button = driver
        .query(By::XPath("//div[@role='button' and contains(.,'mail')]"))
        .single()
        .await?;
    elem_button.click().await?;

    // Find element.
    // let elem_form = driver.find_element(By::Id("search-form")).await?;

    // Find element from element.
    // let elem_text = elem_form.find_element(By::Id("searchInput")).await?;

    // Type in the search terms.
    // elem_text.send_keys("selenium").await?;

    // Always explicitly close the browser. There are no async destructors.
    //driver.quit().await?;

    Ok(())
}
