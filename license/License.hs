{- Haskell code for accessing OSI license API

Copyright © 2016 Clint Adams

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE
-}

{-# LANGUAGE TemplateHaskell #-}

module License (
    OSIIdentifier(..)
  , OSILink(..)
  , OSIOtherName(..)
  , OSIText(..)
  , OSILicense(..)
  , allLicenses
  , licensesMatchingKeyword
  , licenseById
  , licenseBySchemeAndIdentifier
) where

import Control.Monad.Trans.Except (ExceptT(..))
import Data.Aeson (eitherDecode)
import Data.Aeson.TH (defaultOptions, deriveJSON, Options(..))
import Data.Char (toLower)
import Data.Text (Text)
import Network.HTTP.Client (httpLbs, newManager, parseUrl, responseBody)
import Network.HTTP.Client.TLS (tlsManagerSettings)

data OSIIdentifier = OSIIdentifier {
    oiIdentifier :: Text
  , oiScheme :: Text
} deriving (Eq, Read, Show)
$(deriveJSON defaultOptions{fieldLabelModifier = (map toLower . drop 2), constructorTagModifier = map toLower} ''OSIIdentifier)

data OSILink = OSILink {
    olNote :: Text
  , olUrl :: Text
} deriving (Eq, Read, Show)
$(deriveJSON defaultOptions{fieldLabelModifier = (map toLower . drop 2), constructorTagModifier = map toLower} ''OSILink)

data OSIOtherName = OSIOtherName {
    oonName :: Text
  , oonNote :: Maybe Text
} deriving (Eq, Read, Show)
$(deriveJSON defaultOptions{fieldLabelModifier = (map toLower . drop 3), constructorTagModifier = map toLower} ''OSIOtherName)

data OSIText = OSIText {
    otMedia_type :: Text
  , otTitle :: Text
  , otURL :: Text
} deriving (Eq, Read, Show)
$(deriveJSON defaultOptions{fieldLabelModifier = (map toLower . drop 2), constructorTagModifier = map toLower} ''OSIText)

data OSILicense = OSILicense {
    olId :: Text
  , olName :: Text
  , olSuperseded_by :: Maybe Text
  , olKeywords :: [Text]
  , olIdentifiers :: Maybe [OSIIdentifier]  -- BUG in API response?
  , olLinks :: [OSILink]
  , olOther_names :: Maybe [OSIOtherName]   -- BUG in API response?
  , olText :: [OSIText]
} deriving (Eq, Read, Show)
$(deriveJSON defaultOptions{fieldLabelModifier = (map toLower . drop 2), constructorTagModifier = map toLower} ''OSILicense)

getLicenses :: String -> ExceptT String IO [OSILicense]
getLicenses k = ExceptT $ do
  manager <- newManager tlsManagerSettings
  request <- parseUrl $ "https://api.opensource.org/licenses/" ++ k
  response <- httpLbs request manager

  return . eitherDecode . responseBody $ response

allLicenses :: ExceptT String IO [OSILicense]
allLicenses = getLicenses ""

licensesMatchingKeyword :: String -> ExceptT String IO [OSILicense]
licensesMatchingKeyword = getLicenses

getLicense :: String -> ExceptT String IO OSILicense
getLicense k = ExceptT $ do
  manager <- newManager tlsManagerSettings
  request <- parseUrl $ "https://api.opensource.org/license/" ++ k
  response <- httpLbs request manager

  return . eitherDecode . responseBody $ response

licenseById :: String -> ExceptT String IO OSILicense
licenseById = getLicense

licenseBySchemeAndIdentifier :: String -> String -> ExceptT String IO OSILicense
licenseBySchemeAndIdentifier s i = getLicense (concat [s, "/", i])
