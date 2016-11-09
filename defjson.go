package main

import "encoding/xml"

type banqueMat struct{
  Id int64
  Category int64
  Count int64
}

/*type banqueMatXml struct{
  XMLName xml.Name `xml:"Mat"`
  Id int    `xml:"id,attr"`
  Category int `xml:"category"`
  Count int  `xml:"count"`
}*/

type maClef struct{
  XMLName xml.Name `xml:"clef"`
  Id string `xml:"id"`
}


type items struct{
  Id int64
}

type objet struct{
  Id int64
  Chat_link string
  Name string
  Icon string
  Description string
  Type string
  Rarity string
  Level int64
  Vendor_value int64
  Default_skin int64
  Flags string
  GameType string
  Restrictions string
  Details []string
}

type mat struct{
  Name string
  Description string
  Type string
  Level int
  Rarity string
  VendorValue int
  GameType []string
  Flags []string
  Restrictions []string
  Id int
  ChatLink string
  Icon string
}

type armor struct{
  Name string
  Description string
  Type string
  Level int
  Rarity string
  VendorValue int
  GameType []string
  Flags []string
  Restrictions []string
  Id int
  ChatLink string
  Icon string
  Details adetails
}

type adetails struct{
  Type string
  Weight string
  Defense int
  InfusSlot []string
  InfixUpgrade infixUpgrade
  SuffixItim int
  SecondarySuffix string
}

type infixUpgrade struct{
  Id int
  Attributes []string
}


type weapon struct{
  Name string
  Description string
  Type string
  Level int
  Rarity string
  VendorValue int
  GameType []string
  Flags []string
  Restrictions []string
  Id int
  ChatLink string
  Icon string
  Details wdetails
}

type wdetails struct{
  Type string
  DamageType string
  MinPower int
  MaxPower int
  Defense int
  InfusSlot []string
  InfixUpgrade infixUpgrade
  SuffixItim int
  SecondarySuffix string
}

type consumable struct{
  Name string
  Type string
  Level int
  Rarity string
  VendorValue int
  GameType []string
  Flags []string
  Restrictions []string
  ChatLink string
  Icon string
  Details cdetails
}

type cdetails struct{
  Type string
  Duration int64
  ApplyCount int
  Name string
  Icon string
  Description string
}


type price struct{
Id int64
Whitelisted bool
Buys buys
Sells sells
}

type buys struct{
  Quantity int64
  Unit_price int64
}

type sells struct{
  Quantity int64
  Unit_price int64
}
