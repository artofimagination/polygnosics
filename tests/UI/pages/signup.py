"""
This module contains SignupPage,
the page object for the signup page.
"""

from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait


class SignUpPage:
    URL = "http://0.0.0.0:8081/auth_signup"

    USERNAME_FIELD = (By.ID, 'username')
    EMAIL_FIELD = (By.ID, 'email')
    PSW_FIELD = (By.ID, 'psw')
    PSW_REPEAT_FIELD = (By.ID, 'psw-repeat')

    def __init__(self, browser):
        self.browser = browser

    def load(self):
        self.browser.get(self.URL)

    def title(self):
        WebDriverWait(self.browser, 0.1).until(
          EC.presence_of_element_located(self.USERNAME_FIELD))
        return self.browser.title
