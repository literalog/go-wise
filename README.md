<a name="readme-top"></a>



<div align="center">

[![Contributors][contributors-shield]][contributors-url] [![Forks][forks-shield]][forks-url] [![Stargazers][stars-shield]][stars-url] [![Issues][issues-shield]][issues-url] [![AGPL][license-shield]][license-url]

</div>



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/literalog/go-wise">
    <img src="https://i.imgur.com/QD57fTQ.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Wise</h3>

  <p align="center">
    Databases abstraction for Go
    <br />
    <a href="https://github.com/literalog/go-wise/tree/main/docs"><strong>Explore the docs »</strong></a>
    <br />
    <a href="https://github.com/literalog/go-wise/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/literalog/go-wise/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
    ·
    <a href="https://github.com/literalog/go-wise/issues">Open Issues</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#installation">Installation</a>
    </li>
    <li><a href="#usage">Usage</a></li>
      <ol>
        <li><a href="#usage-mongodb">MongoDB</a></li>
      </ol>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Go wise with Wise!

Wise is a Go package offering database abstraction, simplifying database operations and interactions in Go applications. With Wise, developers can seamlessly work with different types of databases without worrying about the underlying implementation details.

This project aims to streamline database handling in Go applications, allowing developers to focus more on building features rather than dealing with the intricacies of database connections and queries.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- Installation -->
### Installation
Please note that this assumes you already have a Go environment set up. If not, you can check the documentation on the official [Go website](https://go.dev/doc/install).

```
go get -u github.com/literalog/go-wise/wise
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Wise supports multiple databases, offering flexibility in selecting the right one for your application.

Below are quick integration guides for supported databases.

<details id="usage-mongodb">

  <summary>MongoDB</summary>

<br/>

Consider you have the following model:
```go
type User struct {
  ID string `json:"id"`
  Name string `json:"name"`
}
```

You can create your repository as follows:

```go
type UserRepository interface {
  wise.MongoRepository[User] 
}

// or

type UserRepository wise.MongoRepository[User]
```

If your document structure differs, you will need to [tailor your own serializer](./docs/mongodb.md##Serializer).

Now, instead of repeatedly crafting a MongoDB implementation, you can simply invoke Wise:

```go
func GetUser(ctx context.Context, id string) (User, error) {
  return UserRepository.Find(ctx, id)
}
```

You can access more in-depth documentation by clicking [here](./docs/mongodb.md).
</details>

<details>
  <summary>Coming soon...</summary>
</details>

<br/>

_For more examples, please refer to the [Documentation](./docs/)_.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the GNU Affero General Public License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/literalog/go-wise.svg?style=for-the-badge
[contributors-url]: https://github.com/literalog/go-wise/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/literalog/go-wise.svg?style=for-the-badge
[forks-url]: https://github.com/literalog/go-wise/network/members
[stars-shield]: https://img.shields.io/github/stars/literalog/go-wise.svg?style=for-the-badge
[stars-url]: https://github.com/literalog/go-wise/stargazers
[issues-shield]: https://img.shields.io/github/issues/literalog/go-wise.svg?style=for-the-badge
[issues-url]: https://github.com/literalog/go-wise/issues
[license-shield]: https://img.shields.io/github/license/literalog/go-wise.svg?style=for-the-badge
[license-url]: https://github.com/literalog/go-wise/blob/main/LICENSE