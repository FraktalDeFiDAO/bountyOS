import os
from dotenv import load_dotenv
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from pymongo import MongoClient
import datetime
import uuid
import time
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC



import config  # Import your config file

# Fetch sensitive details from config
API_KEY = config.API_KEY
USERNAME = config.USERNAME
PASSWORD = config.PASSWORD
MONGO_URI = config.MONGO_URI
CHROMEDRIVER_PATH = config.CHROMEDRIVER_PATH


client = MongoClient(MONGO_URI)
db = client.twitter_trends
collection = db.trends

# Scraper API proxy setup
SCRAPERAPI_PROXY = f"http://api.scraperapi.com:8000?api_key={API_KEY}"

# Set up Chrome options
options = webdriver.ChromeOptions()
options.add_argument(f"--proxy-server={SCRAPERAPI_PROXY}")
options.add_argument("--ignore-certificate-errors")
options.add_argument("--disable-web-security")
options.add_argument("--headless")

# Set up the driver
service = Service(CHROMEDRIVER_PATH)
driver = webdriver.Chrome(service=service, options=options)

driver.get("https://x.com/i/flow/login")

# Login process
username_field = WebDriverWait(driver, 20).until(
    EC.presence_of_element_located((By.CSS_SELECTOR, "input.r-30o5oe.r-1dz5y72.r-13qz1uu.r-1niwhzg.r-17gur6a.r-1yadl64.r-deolkf.r-homxoj.r-poiln3.r-7cikom.r-1ny4l3l.r-t60dpp.r-fdjqy7"))
)

username_field.send_keys(USERNAME)
username_field.send_keys(Keys.RETURN)


password_field = WebDriverWait(driver, 20).until(
    EC.presence_of_element_located((By.CSS_SELECTOR, "input.r-30o5oe.r-1dz5y72.r-13qz1uu.r-1niwhzg.r-17gur6a.r-1yadl64.r-deolkf.r-homxoj.r-poiln3.r-7cikom.r-1ny4l3l.r-t60dpp.r-fdjqy7"))
)
password_field.send_keys(PASSWORD)
password_field.send_keys(Keys.RETURN)

# Wait for trends to load and fetch them
WebDriverWait(driver, 20).until(
    EC.presence_of_element_located((By.XPATH, '//span[contains(text(), "Trending")]'))
)

trends = driver.find_elements(By.XPATH, '//span[contains(text(), "Trending")]/ancestor::div[2]/following-sibling::div[1]//span')

trending_topics = [trend.text for trend in trends if trend.text.strip()]
trending_topics = list(set(trending_topics))

# Print and store the trending topics in MongoDB
for topic in trending_topics:
    print(topic)

unique_id = str(uuid.uuid4())
result = {
    "unique_id": unique_id,
    "trends": trending_topics,
    "date_time": datetime.datetime.now(),
    "ip_address": SCRAPERAPI_PROXY.split('@')[-1]
}
collection.insert_one(result)
print("Trending topics stored in MongoDB.")

# Close the driver
driver.quit()
