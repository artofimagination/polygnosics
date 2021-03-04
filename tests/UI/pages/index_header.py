"""
This module contains IndexPageHeader,
the page object for the header of the main index page.
"""

from selenium.webdriver.common.by import By
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait

class IndexPageHeader:
  URL = "http://polygnosics.localhost"

  SIGNUP_BUTTON = (By.LINK_TEXT, 'Join Us')
  SIGNIN_BUTTON = (By.LINK_TEXT, 'Sign In')
  PRELOADER = (By.CLASS_NAME, 'loader-wrap')

  def __init__(self, browser):
    self.browser = browser

  def load(self):
    self.browser.get(self.URL)

  def maximizeWindow(self):
    self.browser.maximize_window()

  def title(self):
    return self.browser.title

  def goToSignup(self):
    signup = self.browser.find_element(*self.SIGNUP_BUTTON)
    signup.click()

  def goTpSignin(self):
    signin = self.browser.find_element(*self.SIGNIN_BUTTON)
    signin.click()

  def waitUntilPreloaderDisappears(self):
    # wait for loading element to appear
    # - required to prevent prematurely checking if element
    #   has disappeared, before it has had a chance to appear
    WebDriverWait(self.browser, 0.1).until(EC.presence_of_element_located(self.PRELOADER))
    # then wait for the element to disappear
    WebDriverWait(self.browser, 5).until_not(lambda x: x.find_element(*self.PRELOADER).is_displayed())
      
    

