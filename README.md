# phantom[![PkgGoDev](https://pkg.go.dev/badge/github.com/tenntenn/phantom)](https://pkg.go.dev/github.com/tenntenn/phantom)

> [!WARNING]
> **phantom is currently in experimental**. Breaking changes may occur.

`phantom` checks assignabillity of type arguments of phantom type parameter with type alias as folowing codes.

```go
type UUID[T any] = uuid.UUID

type Organization struct {
    OrganizationID UUID[Organization]
    // ...
}

type User struct {
    OrganizationID  UUID[Organization]
    UserID          UUID[User]
    // ...
}

func getUser(orgID UUID[Organization], userID UUID[User]) (*User, error) {
    // ...
}

func run() error {
    var userID UUID[UserID]         // = ...
    var orgID UUID[OrganizationID]  // = ...

    user, err := getUser(userID, orgID) // type checking does not report error but phantom reports error
    // ...
}
```

See [tests](./testdata/src/a/a.go).

## License

This project is licensed under the [MIT License](./LICENSE).

Contributions are always welcome! Feel free to open issues or PRs for bugs and enhancements.
