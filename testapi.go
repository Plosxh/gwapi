package main

import(
  "encoding/json"
  "encoding/xml"
  "net/http"
  "fmt"
  "os"
	"encoding/csv"
  "io/ioutil"
  "database/sql"
  "log"
  //"io"
  _ "github.com/mattn/go-sqlite3"
  "github.com/bwmarrin/discordgo"
//"strings"
	"time"
  "strconv"
  //strings"
)

type item struct{
  XMLName xml.Name `xml:"items"`
  BanqueMatXml []banqueMatXml `xml:"mat"`
}
type banqueMatXml struct{
  XMLName xml.Name `xml:"mat"`
  Id int64    `xml:"id"`
  Category int64 `xml:"category"`
  Count int64  `xml:"count"`
}

var (
	BotID    string
)

func main() {

  //https://api.guildwars2.com/v2/tokeninfo?access_token=65D84368-DA6E-9D4A-8B6E-70C0395432961B8D9A2D-1F1E-4F28-B484-9D0DFE20DBFF
   //clef := "65D84368-DA6E-9D4A-8B6E-70C0395432961B8D9A2D-1F1E-4F28-B484-9D0DFE20DBFF"
   var choix int
   var id int64
   var category int
   var count int
   var objet string
   var item_id int
   var name string
   //var mesObjets string
   var nb int
   var Ids []int64
   var monProfit int
   var fin string
   var monTime int
   //var loop bool
   fmt.Println("Choisissez : 1-Mettre à jour la Banque, 2-Prix de chaque item en banque, 3-halloween mode :")
   _,err := fmt.Scanln(&choix)
   if err != nil {
     log.Fatal(err)
   }


  switch choix {
  case 1:
  checkBank(getClef())

  case 5:
  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM Bank")
    for rows.Next(){
    err = rows.Scan(&id,&name,&item_id,&category,&count)
    fmt.Println("ID : ", item_id," Nom : ",name, " Category : ",category," Count : ",count)
    if err != nil {
      log.Fatal(err)
    }
  }
  fmt.Println("Choisissez l'Id d'un objet : ")
  _,err = fmt.Scanln(&objet)

  getUnItem(objet)

  case 3:
  //for {
  x := time.Now()
  halloween(x)
  doEveryhalloween(300*time.Second)

  //}
  /*loop = true
  for  loop {
    fmt.Println("Choisissez l'Id d'un objet ou stop pour arrêter : ")
    _,err = fmt.Scanln(&objet)
    if objet =="stop"{
      loop=false
    }else{
      loop=true
      mesObjets=objet

      getUnItem(mesObjets)
    }

  }*/

  case 4:
  addFav()

case 2:
  fmt.Println("Choisissez un profit minimum (entrez un entier entre 0 et 100) : ")
  _,err = fmt.Scanln(&monProfit)

  fmt.Println("Choisissez le temps de rafraichissement en second : ")
  _,err = fmt.Scanln(&monTime)
  nb=0
  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM Bank")
    for rows.Next(){
    err = rows.Scan(&id,&item_id,&name,&category,&count)
    //fmt.Println("ID : ", item_id," Nom : ",name, " Category : ",category," Count : ",count)
    Ids=append(Ids,int64(item_id))
    if err != nil {
      log.Fatal(err)
      fmt.Println("arret scanln in case 5")
    }
    nb++
  }
  x := time.Now()
  getBankPrices(x,Ids,monProfit)
  doEvery(300*time.Second,monProfit,Ids)
  //getBankPrices(Ids,monProfit)
  case 6:
    getInvendable()
  case 7:
    checkFav()
  case 8:
    dg, err := discordgo.New("pierre.charrat@etu.univ-lyon1.fr", "gwapipass", "MjQwNzYxMDQ5Njk4MDA5MDg4.CvIBxg.ylEFpmzWJ1oVsq2lDylH9pABkC8")
    if err != nil {
      fmt.Println("error creating Discord session,", err)
      return
    }

    u, err := dg.User("@me")
    if err != nil {
      fmt.Println("error obtaining account details,", err)
    }

    // Store the account ID for later use.
    BotID = u.ID

    dg.AddHandler(messageCreate)

    // Open the websocket and begin listening.
    err = dg.Open()
    if err != nil {
      fmt.Println("error opening connection,", err)
      return
    }

    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    // Simple way to keep program running until CTRL-C is pressed.
    <-make(chan struct{})
    return
  case 9:
    toxml()

  }
  //doEvery(10*time.Second)
  //mesItems:=getItems()

  //fmt.Println(mesItems[0])


  /*
    foo2 := price{}
    getJson("https://api.guildwars2.com/v2/commerce/prices?id=19684", &foo2)
    fmt.Println(foo2.Buys.UnitePrice)*/
    fmt.Println("appuyer sur entrer pour fermer : ")
    _,err = fmt.Scanln(&fin)

}


func halloween(t time.Time){
  url := "https://api.guildwars2.com/v2/commerce/prices?ids=36038,67379,67386,36041,36074,36081,36080,36077,36084,36076,67367,67371,67368,36095,36060,48807,36061,48806,36059,48805,72113,67380,36103,67369,67381,36047,36066,67382,36102,36050,36065,79637,79638,71931,71946,70732,79690,76131,76642,67370,67372,67375,79647"
  var mesPrices []price
  getJson(url,&mesPrices)

  for i := 0; i < len(mesPrices); i++ {
    //fmt.Println(I)
    if mesPrices[i].Buys.Unit_price == int64(0) && mesPrices[i].Sells.Unit_price == int64(0){
        supprItem(mesPrices[i].Id)
        fmt.Println("Item non vendable : ",mesPrices[i].Id)
  }else{
    profit :=calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)

    if profit>=float64(0){
      fmt.Println("Nom : ",getNom(mesPrices[i].Id)," | Achat : ",mesPrices[i].Buys.Unit_price," | Vente : ",mesPrices[i].Sells.Unit_price-1," | Profit : ",profit,"%")
      fmt.Println("--------------------------------------------------------------------------------------------------------")
    }else{
      //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
      }
    }
  }
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  var id int64
  var category int
  var count int
  var item_id int
  var name string
  var Ids []int64
  //var out1 string
  //var out2 string
  //var out3 string
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}
  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM Bank")
    for rows.Next(){
    err = rows.Scan(&id,&item_id,&name,&category,&count)
    //fmt.Println("ID : ", item_id," Nom : ",name, " Category : ",category," Count : ",count)
    Ids=append(Ids,int64(item_id))
    if err != nil {
      log.Fatal(err)
      fmt.Println("arret scanln in case 5")
    }
  }

  mesPrices,mesPrices2,mesPrices3:=getDiscordBankPrices(time.Now(),Ids,20)
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!prices" {
    //  mesItems:=getItems()
    for i := 0; i < len(mesPrices); i++ {
      profit :=calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)
      if profit>=float64(50){
        //out1 += "Nom : "+getNom(mesPrices[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices[i].Sells.Unit_price-1,10)+" | Profit : "+strconv.FormatInt(int64(profit),10)+"%"+"\n"+"--------------------------------------------------------------------------------------------------------"+"\n"
    _, _ = s.ChannelMessageSend(m.ChannelID, "Nom : "+getNom(mesPrices[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices[i].Sells.Unit_price-1,10)+"\n"+"--------------------------------------------------------------------------------------------------------")
    //_, _ = s.ChannelMessageSend(m.ChannelID,"--------------------------------------------------------------------------------------------------------")

      }
    }
  //  _, _ = s.ChannelMessageSend(m.ChannelID,out1)
    for i := 0; i < len(mesPrices2); i++ {
      profit :=calcFees(mesPrices3[i].Buys.Unit_price,mesPrices3[i].Sells.Unit_price)
      if profit>=float64(50){
        //out2 += "Nom : "+getNom(mesPrices2[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices2[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices2[i].Sells.Unit_price-1,10)+" | Profit : "+strconv.FormatInt(int64(profit),10)+"%"+"\n"+"--------------------------------------------------------------------------------------------------------"+"\n"
    _, _ = s.ChannelMessageSend(m.ChannelID, "Nom : "+getNom(mesPrices2[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices2[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices2[i].Sells.Unit_price-1,10)+"\n"+"--------------------------------------------------------------------------------------------------------")
    //_, _ = s.ChannelMessageSend(m.ChannelID,"--------------------------------------------------------------------------------------------------------")
      }
    }
    //_, _ = s.ChannelMessageSend(m.ChannelID,out2)

    for i := 0; i < len(mesPrices3); i++ {
      profit :=calcFees(mesPrices3[i].Buys.Unit_price,mesPrices3[i].Sells.Unit_price)
      if profit>=float64(50){
    _,_ = s.ChannelMessageSend(m.ChannelID, "Nom : "+getNom(mesPrices3[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices3[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices3[i].Sells.Unit_price-1,10)+"\n"+"--------------------------------------------------------------------------------------------------------")
    //_, _ = s.ChannelMessageSend(m.ChannelID,"--------------------------------------------------------------------------------------------------------")
      //out3 += "Nom : "+getNom(mesPrices3[i].Id)+" | Achat : "+strconv.FormatInt(mesPrices3[i].Buys.Unit_price,10)+" | Vente : "+strconv.FormatInt(mesPrices3[i].Sells.Unit_price-1,10)+" | Profit : "+strconv.FormatInt(int64(profit),10)+"%"+"\n"+"--------------------------------------------------------------------------------------------------------"+"\n"
      }
    }
    //_, _ = s.ChannelMessageSend(m.ChannelID,out3)

    _,_ = s.ChannelMessageSend(m.ChannelID,"Fin")
    //fmt.Println(string(mesItems[0].Id))
	}

	// If the message is "pong" reply with "Ping!"
}

func toxml(){
  //var mesObjets string
  var Ids []int64
  getJson("https://api.guildwars2.com/v2/items",&Ids)
  for i:=25464; i < len(Ids);i++{
    var monObjet objet
    getJson("https://api.guildwars2.com/v2/items/"+strconv.FormatInt(Ids[i],10),&monObjet)

    //fmt.Println(i)
    output,err :=json.MarshalIndent(&monObjet, "", "\t\t")
      if err != nil {
        fmt.Printf("error: %v\n", err)
      }
    /*  if i == 0{
        mesObjets += "["+string(output)
      }else{
        mesObjets += ","+ string(output)
      }*/
      errecri := ioutil.WriteFile("./items/"+monObjet.Type+"/"+strconv.Itoa(i)+".json",output,0644)
      if errecri != nil {
        fmt.Printf("error: %v\n", errecri)
      }

  }
  /*mesObjets += "]"
  err := ioutil.WriteFile("out.json", []byte(mesObjets), 0644)
  if err != nil {
    fmt.Printf("error: %v\n", err)
  }*/
}


func getInvendable(){
  var monItem price
  var id int64
  var item_id int64
  var name string
  var category int64
  var count int64
  var Ids []int64
  url :="https://api.guildwars2.com/v2/commerce/prices?id="
  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM Bank")
    for rows.Next(){
    err = rows.Scan(&id,&item_id,&name,&category,&count)
    //fmt.Println("ID : ", item_id," Nom : ",name, " Category : ",category," Count : ",count)
    Ids=append(Ids,int64(item_id))
    if err != nil {
      log.Fatal(err)
      fmt.Println("arret scanln in case 5")
    }

    for i := 0; i < len(Ids); i++ {
      getJson(url+strconv.FormatInt(Ids[i],10),&monItem)
      if monItem.Buys.Unit_price ==0{
        if monItem.Sells.Unit_price ==0{
          fmt.Println("Cet Item est invendable :",monItem.Id)
        }
      }

    }
  }

}


func getUnItem(I string)  price{

  url := "https://api.guildwars2.com/v2/commerce/prices?id="+I
  //fmt.Println(url)

  var Unitems price
  getJson(url,&Unitems)
  fmt.Println("item: ",getNom(Unitems.Id))
  //for i := 0; i < len(Unitems); i++ {
      fmt.Println("Achat : ",Unitems.Buys.Unit_price," Vente : ",Unitems.Sells.Unit_price," Profit : ",calcFees(Unitems.Buys.Unit_price,Unitems.Sells.Unit_price))

  if Unitems.Buys.Unit_price != 0 || Unitems.Sells.Unit_price != 0{
        addCsvTest(Unitems,getNom(Unitems.Id))
  }
  //}
  return Unitems
}


func addFav()  {

  var choix int64
  var id int
  var name string
  var item_id int
  var category int
  var count int

  //var id2 int
  var name2 string
  var item_id2 int
  var category2 int
  var count2 int
  var expr string

  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM Bank")
    for rows.Next(){
    err = rows.Scan(&id,&item_id,&name,&category,&count)
    fmt.Println("ID : ", item_id," Nom : ",name, " Category : ",category," Count : ",count)
    if err != nil {
      log.Fatal(err)
    }
  }

  fmt.Println("Choisissez l'id d'un item à ajouter en favori : ")
  _,err = fmt.Scanln(&choix)

  //rows,err =db.Query("SELECT * FROM Bank where id="+choix)
  //err = rows.Scan(&id2,&item_id2,&name2,&category2,&count2)
    fmt.Println("ID : ", item_id2," Nom : ",name2, " Category : ",category2," Count : ",count2)
    expr ="INSERT INTO favori VALUES ("+strconv.FormatInt(choix,10) +",\""+getNom(choix)+"\")"
    fmt.Println(expr)
    _,err =db.Exec(expr)
  if err != nil {
    log.Fatal(err)
    }


}


func checkFav(){
  var mesPrices []price
  var mesPrices2 []price
  var mesPrices3 []price
  var profit float64
  var Ids []int64
  var item_id int64
  var name string
  var min int64

  fmt.Println("Choisissez la valeur de profit minimum ( en pourcent) : ")
  _,err:= fmt.Scanln(&min)

  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    rows,err :=db.Query("SELECT * FROM favori")
    for rows.Next(){
    err = rows.Scan(&item_id,&name)
    Ids=append(Ids,int64(item_id))
    if err != nil {
      log.Fatal(err)
    }
  }


  //L'api est limité à 200 items à la fois du coup on sépare les 413 items en 3
  //fmt.Println("len de I : ",len(I))
  if len(Ids)<200 {
    url := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[0],10)
    for i := 1; i < len(Ids); i++ {
      url = url +","+strconv.FormatInt(Ids[i],10)
    }
    getJson(url,&mesPrices)

    for i := 0; i < len(mesPrices); i++ {
      //fmt.Println(I)
      if mesPrices[i].Buys.Unit_price == int64(0) && mesPrices[i].Sells.Unit_price == int64(0){
            supprItem(mesPrices[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices[i].Id)," | Achat : ",mesPrices[i].Buys.Unit_price," | Vente : ",mesPrices[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }
    fmt.Println("============================================================================================")
  }

  if len(Ids)>200 {
    url := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[0],10)
    for i := 1; i < 199; i++ {
      url = url +","+strconv.FormatInt(Ids[i],10)
    }

    url2 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[200],10)
    for i := 201; i < len(Ids); i++ {
      url2 = url2 +","+strconv.FormatInt(Ids[i],10)

    }

    getJson(url,&mesPrices)
    getJson(url2,&mesPrices2)
    for i := 0; i < len(mesPrices); i++ {
      //fmt.Println(I)
      if mesPrices[i].Buys.Unit_price == int64(0) && mesPrices[i].Sells.Unit_price == int64(0){
            supprItem(mesPrices[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices[i].Id)," | Achat : ",mesPrices[i].Buys.Unit_price," | Vente : ",mesPrices[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }


    for i := 0; i < len(mesPrices2); i++ {
      //fmt.Println(I)
      if mesPrices2[i].Buys.Unit_price == 0 && mesPrices2[i].Sells.Unit_price == 0{
            supprItem(mesPrices2[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices2[i].Buys.Unit_price,mesPrices2[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices2[i].Id)," | Achat : ",mesPrices2[i].Buys.Unit_price," | Vente : ",mesPrices2[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }
    fmt.Println("============================================================================================")
  }


  if len(Ids)>400 {
    url := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[0],10)
    for i := 1; i < 199; i++ {
      url = url +","+strconv.FormatInt(Ids[i],10)
    }

    url2 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[200],10)
    for i := 201; i < 399; i++ {
      url2 = url2 +","+strconv.FormatInt(Ids[i],10)

    }

    url3 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(Ids[400],10)
    for i := 401; i < len(Ids); i++ {
      url3 = url3 +","+strconv.FormatInt(Ids[i],10)

    }
    getJson(url,&mesPrices)
    getJson(url2,&mesPrices2)
    getJson(url3,&mesPrices3)

    for i := 0; i < len(mesPrices); i++ {
      //fmt.Println(I)
      if mesPrices[i].Buys.Unit_price == int64(0) && mesPrices[i].Sells.Unit_price == int64(0){
            supprItem(mesPrices[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices[i].Id)," | Achat : ",mesPrices[i].Buys.Unit_price," | Vente : ",mesPrices[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }


    for i := 0; i < len(mesPrices2); i++ {
      //fmt.Println(I)
      if mesPrices2[i].Buys.Unit_price == 0 && mesPrices2[i].Sells.Unit_price == 0{
            supprItem(mesPrices2[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices2[i].Buys.Unit_price,mesPrices2[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices2[i].Id)," | Achat : ",mesPrices2[i].Buys.Unit_price," | Vente : ",mesPrices2[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }


    for i := 0; i < len(mesPrices3); i++ {
      //fmt.Println(I)
      if mesPrices3[i].Buys.Unit_price == 0 && mesPrices3[i].Sells.Unit_price == 0{
            supprItem(mesPrices3[i].Id)
            fmt.Println("Item non vendable : ",mesPrices[i].Id)
      }else{
        profit =calcFees(mesPrices3[i].Buys.Unit_price,mesPrices3[i].Sells.Unit_price)

        if profit>=float64(min){
          fmt.Println("Nom : ",getNom(mesPrices3[i].Id)," | Achat : ",mesPrices3[i].Buys.Unit_price," | Vente : ",mesPrices3[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        }else{
          //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
        }
      }
    }
  fmt.Println("============================================================================================")
  }

  //fmt.Println(url)

  //fmt.Println(len(mesPrices))


}

func getBankPrices(t time.Time,I []int64,min int)  {
  var mesPrices []price
  var mesPrices2 []price
  var mesPrices3 []price
  var profit float64
  fmt.Println(t)
  //L'api est limité à 200 items à la fois du coup on sépare les 413 items en 3
  //fmt.Println("len de I : ",len(I))
  url := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[0],10)
  for i := 1; i < 199; i++ {
    url = url +","+strconv.FormatInt(I[i],10)
  }

  url2 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[200],10)
  for i := 201; i < 399; i++ {
    url2 = url2 +","+strconv.FormatInt(I[i],10)

  }

  url3 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[400],10)
  for i := 401; i < 413; i++ {
    url3 = url3 +","+strconv.FormatInt(I[i],10)

  }
  //fmt.Println(url)
  getJson(url,&mesPrices)
  getJson(url2,&mesPrices2)
  getJson(url3,&mesPrices3)
  //fmt.Println(len(mesPrices))
  for i := 0; i < len(mesPrices); i++ {
    //fmt.Println(I)
    if mesPrices[i].Buys.Unit_price == int64(0) && mesPrices[i].Sells.Unit_price == int64(0){
          supprItem(mesPrices[i].Id)
          fmt.Println("Item non vendable : ",mesPrices[i].Id)
    }else{
      profit =calcFees(mesPrices[i].Buys.Unit_price,mesPrices[i].Sells.Unit_price)

      if profit>=float64(min){
        fmt.Println("Nom : ",getNom(mesPrices[i].Id)," | Achat : ",mesPrices[i].Buys.Unit_price," | Vente : ",mesPrices[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        fmt.Println("--------------------------------------------------------------------------------------------------------")
      }else{
        //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
      }
    }
  }


  for i := 0; i < len(mesPrices2); i++ {
    //fmt.Println(I)
    if mesPrices2[i].Buys.Unit_price == 0 && mesPrices2[i].Sells.Unit_price == 0{
          supprItem(mesPrices2[i].Id)
          fmt.Println("Item non vendable : ",mesPrices[i].Id)
    }else{
      profit =calcFees(mesPrices2[i].Buys.Unit_price,mesPrices2[i].Sells.Unit_price)

      if profit>=float64(min){
        fmt.Println("Nom : ",getNom(mesPrices2[i].Id)," | Achat : ",mesPrices2[i].Buys.Unit_price," | Vente : ",mesPrices2[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        fmt.Println("--------------------------------------------------------------------------------------------------------")
      }else{
        //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
      }
    }
  }


  for i := 0; i < len(mesPrices3); i++ {
    //fmt.Println(I)
    if mesPrices3[i].Buys.Unit_price == 0 && mesPrices3[i].Sells.Unit_price == 0{
          supprItem(mesPrices3[i].Id)
          fmt.Println("Item non vendable : ",mesPrices[i].Id)
    }else{
      profit =calcFees(mesPrices3[i].Buys.Unit_price,mesPrices3[i].Sells.Unit_price)

      if profit>=float64(min){
        fmt.Println("Nom : ",getNom(mesPrices3[i].Id)," | Achat : ",mesPrices3[i].Buys.Unit_price," | Vente : ",mesPrices3[i].Sells.Unit_price-1," | Profit : ",profit,"%")
        fmt.Println("--------------------------------------------------------------------------------------------------------")
      }else{
        //fmt.Println("Nom : ",getNom(mesPrices[i].Id),"a un profit de : ",profit," ce qui est trop faible.")
      }
    }
  }
  fmt.Println("============================================================================================")
}

func getDiscordBankPrices(t time.Time,I []int64,min int)  ([]price,[]price,[]price){
  var mesPrices []price
  var mesPrices2 []price
  var mesPrices3 []price
  fmt.Println(t)
  //L'api est limité à 200 items à la fois du coup on sépare les 413 items en 3
  //fmt.Println("len de I : ",len(I))
  url := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[0],10)
  for i := 1; i < 199; i++ {
    url = url +","+strconv.FormatInt(I[i],10)
  }

  url2 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[200],10)
  for i := 201; i < 399; i++ {
    url2 = url2 +","+strconv.FormatInt(I[i],10)

  }

  url3 := "https://api.guildwars2.com/v2/commerce/prices?ids="+strconv.FormatInt(I[400],10)
  for i := 401; i < 413; i++ {
    url3 = url3 +","+strconv.FormatInt(I[i],10)

  }
  //fmt.Println(url)
  getJson(url,&mesPrices)
  getJson(url2,&mesPrices2)
  getJson(url3,&mesPrices3)
  //fmt.Println(len(mesPrices))
  return mesPrices,mesPrices2,mesPrices3
}


func supprItem(Id int64){
  db, err := sql.Open("sqlite3", "./itemgw.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    _,err =db.Exec("DELETE from bank where id="+strconv.FormatInt(Id,10))

}

func checkBank(key string)  {
  //var objets items
  var foo1 []banqueMatXml
  //var tempo1 []banqueMatXml
  //var foo2 banqueMatXml
  //fmt.Println("allo ?")
  getJson("https://api.guildwars2.com/v2/account/materials?access_token="+key, &foo1)
  fmt.Println("vous avez : ",len(foo1)," objects dans vos materiaux.")
    //fmt.Println("allo 2 ?")

  writer,_ :=os.OpenFile("./gwitem.xml", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  enc := xml.NewEncoder(writer)
  //fmt.Println("test")
  //enc.Indent("  ", "    ")
  //fmt.Println(len(foo1))
  /*for i := 0; i < len(foo1); i++ {
  fmt.Println("bla")

  tempo1[i].Id = foo1[i].Id
  tempo1[i].Category = foo1[i].Category
  tempo1[i].Count = foo1[i].Count
  }*/
  foo2 := &item{BanqueMatXml:foo1}
  //objets = append(objets,foo1[i].Id,foo1[i].Category,foo1[i].Count)

  /*foo2.Id = foo1[i].Id
  foo2.Category = foo1[i].Category
  foo2.Count = foo1[i].Count*/
  //fmt.Println(foo2)
  if err := enc.Encode(foo2); err != nil {
      fmt.Printf("error: %v\n", err)
    }



  var monItem item
  xmlContent, _ := ioutil.ReadFile("gwitem.xml")
  err := xml.Unmarshal(xmlContent, &monItem)
  //err = xml.Unmarshal(xmlContent, &R)
  //fmt.Println(monItem)
  if err != nil { panic(err) }
  itemlen := len(monItem.BanqueMatXml)

    db, err := sql.Open("sqlite3", "./itemgw.db")
    if err != nil {
      log.Fatal(err)
    }
    defer db.Close()

    for i := 0; i < itemlen; i++ {
      var name = getNom(monItem.BanqueMatXml[i].Id)
      _,err =db.Exec("INSERT INTO bank VALUES (NULL,"+strconv.FormatInt(monItem.BanqueMatXml[i].Id,10)+",\""+name+"\","+strconv.FormatInt(monItem.BanqueMatXml[i].Category,10)+","+strconv.FormatInt(monItem.BanqueMatXml[i].Count,10)+")")


    }
}

func getClef()string{

  var clef maClef
  xmlContent, _ := ioutil.ReadFile("apikey.xml")
  err := xml.Unmarshal(xmlContent, &clef)
  //err = xml.Unmarshal(xmlContent, &R)
  if err != nil { panic(err) }
  //fmt.Println(primlen)
  return clef.Id
}


func getItems()  []items{
  url := "https://api.guildwars2.com/v2/items"

  var mesItems []items

  getJson(url,&mesItems)
  return mesItems

}


func doEvery(d time.Duration,p int, i []int64) {
	for x := range time.Tick(d) {
    getBankPrices(x,i,p)
	}
}

func doEveryhalloween(d time.Duration) {
	for x := range time.Tick(d) {
    halloween(x)
	}
}

func pingApi(t time.Time) price{
  url := "./prices.json"
  var foo1 price // or &Foo{}
  //  fmt.Println(t.Clock)
  getJson("https://api.guildwars2.com/v2/commerce/prices?id=19684", &foo1)
  //getJson(url,foo1)
  //println(foo1.Buys.Unit_price)
  //fmt.Println(foo1)


  file, err := ioutil.ReadFile(url)
    if err != nil {
      fmt.Printf("File error: %v\n", err)
      os.Exit(1)
  }
  //fmt.Printf("%s\n", string(file))

 jsontype := new(price)
  json.Unmarshal(file, &jsontype)
  fmt.Printf("Results: %v\n", jsontype.Id)

  //addCsv(*foo1)
  /*d := json.NewDecoder(strings.NewReader(jsontype))
  d.UseNumber()
  var x interface{}
  if err := d.Decode(&x); err != nil {
      log.Fatal(err)
  }
  fmt.Printf("decoded to %#v\n", x)*/
  return foo1
}

func addCsv(p price,name string) {
  //values := []string{}
  //result := []string{}
  fichier := "./items/"+name +".csv"
  fmt.Println(name)
	f, err := os.OpenFile(fichier, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
    //f, err := os.Create(fichier)
		f,_ := os.Create(fichier)
    f,_ = os.OpenFile(fichier, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    w := csv.NewWriter(f)
    w.Write([]string{strconv.FormatInt(p.Id,10)+";" +getNom(p.Id)+";" +strconv.FormatInt(p.Buys.Unit_price,10)+";" + strconv.FormatInt(p.Sells.Unit_price,10)+";"+strconv.FormatFloat(calcFees(p.Buys.Unit_price,p.Sells.Unit_price),'G',0,64)})
    fmt.Println("erreur")
    w.Flush()
  }else{
    w := csv.NewWriter(f)
    w.Write([]string{strconv.FormatInt(p.Id,10)+";" +getNom(p.Id)+";" +strconv.FormatInt(p.Buys.Unit_price,10)+";" + strconv.FormatInt(p.Sells.Unit_price,10)+";"+strconv.FormatFloat(calcFees(p.Buys.Unit_price,p.Sells.Unit_price),'G',0,64)})
    fmt.Println("pas d'erreur")
    w.Flush()
  }

	//for i := 0; i < 10; i++ {
  /*values[0]=strings.Join(strconv.FormatInt(p.Id,10),";")
  values[1]=strings.Join(getNom(p.Id),";")
  values[2]=strings.Join(strconv.FormatInt(p.Buys.Unit_price,10),";")
  values[3]=strings.Join(strconv.FormatInt(p.Sells.Unit_price,10),";")*/


  //result[0] =strings.Join(values,";")
  //w.Write(values[0],values[1],values[2],values[3])
  //}

}



func addCsvTest(p price,name string) {
  fichier := "./items/"+name +".csv"
  fmt.Println(name)
	f, err := os.OpenFile(fichier, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		f,_ := os.Create(fichier)
    f,_ = os.OpenFile(fichier, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    w := csv.NewWriter(f)
    w.Write([]string{"Id;Name;Buy Price;Sell Price;Profit"})
    w.Flush()
  //  x := csv.NewWriter(f)
    //x.Write([]string{strconv.FormatInt(p.Id,10)+";" +getNom(p.Id)+";" +strconv.FormatInt(p.Buys.Unit_price,10)+";" + strconv.FormatInt(p.Sells.Unit_price,10)+";"+strconv.FormatFloat(calcFees(p.Buys.Unit_price,p.Sells.Unit_price),'G',0,64)})
    fmt.Println("erreur")
    //x.Flush()
  }else{
    w := csv.NewWriter(f)
    w.Write([]string{strconv.FormatInt(p.Id,10)+";" +getNom(p.Id)+";" +strconv.FormatInt(p.Buys.Unit_price,10)+";" + strconv.FormatInt(p.Sells.Unit_price,10)+";"+strconv.FormatFloat(calcFees(p.Buys.Unit_price,p.Sells.Unit_price),'G',0,64)})
    fmt.Println("pas d'erreur")
    w.Flush()
  }
}


func getJson(url string, target interface{}) error {
    r, err := http.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}


func calcFees(buy int64, sell int64) float64{
  var profit float64
  //profit = ((float64(sell)*0.85)-float64(buy))
  if sell == 0 {
    if buy ==0{
      fmt.Println("cette objet n'existe pas en vente !")
    }else{
          profit =0
    }

  }else{
    //profit =(float64(100)-(float64((buy*100.0))/(float64(sell)*0.85)))
    profit =(((float64(sell-1)*0.85)-float64(buy))*float64(100))/float64(buy)
  }

  return profit
}

func getNom(id int64)string{
  var monMat []mat
  getJson("https://api.guildwars2.com/v2/items?ids="+strconv.FormatInt(id,10), &monMat)

  return monMat[0].Name
}
