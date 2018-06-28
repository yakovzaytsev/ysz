# ysz (go)

My go code

- Serve some.html:

        ysz.HandleHomePage()

    HOME_PAGE environment variable should point to some.html

    Example:
    
        srv = &http.Server{Addr: ":8000"}
        ysz.HandleHomePage()


