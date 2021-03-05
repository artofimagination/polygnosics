"""
This module contains SigninPage,
the page object for the signup page.
"""

from selenium.webdriver.common.by import By


class SigninPage:
    URL = "http://0.0.0.0:8081/auth_login"

    EMAIL_FIELD = (By.ID, 'email')
    PSW_FIELD = (By.ID, 'psw')

    def __init__(self, browser):
        self.browser = browser

    def title(self):
        return self.browser.title

    def load(self):
        self.browser.get(self.URL)
