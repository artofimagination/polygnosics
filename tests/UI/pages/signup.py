"""
This module contains SignupPage,
the page object for the signup page.
"""

from selenium.webdriver.common.by import By

class SignUpPage:
  USERNAME_FIELD = (By.ID, 'username')
  EMAIL_FIELD = (By.ID, 'email')
  PSW_FIELD = (By.ID, 'psw')
  PSW_REPEAT_FIELD = (By.ID, 'psw-repeat')

  def __init__(self, browser):
    self.browser = browser

  def title(self):
    self.browser.title