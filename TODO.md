# Todos
- [ ] Create Database


# Notes
The code is used from https://medium.com/geekculture/authenticate-go-graphql-with-jwt-436c74340d

mutation{
  auth{
register(input: {
name: "G",
email: "g@g.g",
password: "12345"
})
  }
}

mutation {
auth{
login(email: "g@g.g", password: "12345")
}
}
