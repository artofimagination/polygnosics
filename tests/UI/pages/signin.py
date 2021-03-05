"""
This module contains SigninPage,
the page object for the signup page.
"""

from selenium.webdriver.common.by import By


class SignUpPage:
    EMAIL_FIELD = (By.ID, 'email')
    PSW_FIELD = (By.ID, 'psw')

    def __init__(self, browser):
        self.browser = browser

    def title(self):
        self.browser.title
