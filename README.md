VENT
======
`A Vent Social App`


1.框架图

        .---------.             .---------.
        | clientA |             | clientB |
        |---------|             |---------|
             |                       |
             V                       V
    .------------------------------------------.                   .---------------------------.
    |                IRIS                      |                   |    consul                 |
    |           http restful api               |------------------>|   use for global config &&|
    |                                          |                   |   service discovery       |
    |                                          |                   |                           |
    .------------------------------------------.                   .---------------------------.
             |                     |                                          |             |
             |                     |        GET SERVICES && LOADBALANCE       |             |
             |       gRPC          | <----------------------------------------.             |
             |                     |                                                        |
             .---------------------.                                                        |
                       ↑                                                                   |
                       |                                                                    |
                       |                                                                    |
                       ↓                                                                   |
    .-----------------------------------------.                                             |
    |    single micro service ,eg:            |                                             |
    |      UserService/relationService/       |        SERVICE REGISTER                     |
    |      authService /captchaService        |---------------------------------------------.
    |    and so on.                           |
    |                                         |
    .-----------------------------------------.
                        ↑
                        |
                        ↓
                   .----------.
                   |   redis  |
                   |          |
                   .----------.
                   
                   
2. TODO
  2.1 Use message broker to do async event
