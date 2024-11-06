# Wallet app idea

Wallet app is a demonstration of my skills in Golang.

Apllication Architecture consist of best practicies of language itself but not limited by them  

Clean Architecture + RESTAPI 'like structure allows to anyone that has ever worked on 
web applications understand project in the matter of seconds.


# Build && Run 
 
Requires zero prep and can be achieved via 

```make
make go
```

# Configuration 

Application uses only Environment Variables for its configuration.

Comes with explicit default configuration 

Sources cosists of exactly none magic-defaults   

# Maintain

Application logs are written in stderr and stdout ( configurable ) in json format
so you can pipe it to any aggregation server or inspect raw 

