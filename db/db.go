package db

type customer struct {
  ID int64
  Name string
  Address string
  Phone string
}

var Customers = []customer{
  customer{295415523, "Виктор Васильевич Кузнецов", "улица Ленина", "+79284517471"},
  customer{1201236817, "Анастасия Викторовна Кузнецова", "улица Ленина 288", "+79284517424"},
}
