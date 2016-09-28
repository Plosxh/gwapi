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
  _ "github.com/mattn/go-sqlite3"
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
   var mesObjets string
   var nb int
   var Ids []int64
   var monProfit int
   var fin string
   var monTime int
   var loop bool
   fmt.Println("Choisissez : 1-Mettre à jour la Banque, 2-Voir les prix, 3-Tester getUnItem, 4-par ajouter un favori, 5-getBankPrice :")
   _,err := fmt.Scanln(&choix)
   if err != nil {
     log.Fatal(err)
   }


  switch choix {
  case 1:
  checkBank(getClef())

  case 2:
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



  //}
  loop = true
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

  }

  case 4:
  addFav()

  case 5:
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
  doEvery(10*time.Second,monProfit,Ids)
  //getBankPrices(Ids,monProfit)
  case 6:
    getInvendable()
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


func getUnItem(I string)  {

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

}


func addFav()  {

  var choix string
  var id int
  var name string
  var item_id int
  var category int
  var count int

  var id2 int
  var name2 string
  var item_id2 int
  var category2 int
  var count2 int

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

  rows,err =db.Query("SELECT * FROM Bank where id="+choix)
  for rows.Next(){
  err = rows.Scan(&id2,&item_id2,&name2,&category2,&count2)
    fmt.Println("ID : ", item_id2," Nom : ",name2, " Category : ",category2," Count : ",count2)
    //_,err =db.Exec("INSERT INTO favori VALUES (\""+choix +"\",\""+name2+"\")")
  if err != nil {
    log.Fatal(err)
    }
  }


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
