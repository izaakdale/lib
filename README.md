# lib
Collection of little helper packages I have curated to ease my development process.

Import:

    go get github.com/izaakdale/lib
    
Packages:

------------------------------------------------

      Name - server
      Description - http.Server setup
      Functions - 
        New - Returns a Server with the specified options from below
        WithHost
        WithPort
        WithTimeouts
        WithTimeoutHandler
        (*Server) Run - Shorthand for ListenAndServe
------------------------------------------------

      Name - router
      Description - http.Handler setup, 
      Functions - 
        New - Returns a Handler with the specified options from below,
              Includes default /_/ping route and a middleware that logs request endpoint and status.
        WithRoute
        WithMiddleware
------------------------------------------------

      Name - listener
      Description - AWS SQS Client
      Functions - 
        Initialise
        Listen
        WithEndpoint
        WithMaxNumerOfMessages
        WithVisibilityTimeout
        WithWaitTimeSeconds
------------------------------------------------

      Name - publisher
      Description - AWS SNS Client
      Functions - 
        Initialise
        Publish
        WithEndpoint
        WithPublisher
------------------------------------------------

      Name - response
      Description - http writer
      Functions - 
        WriteJson
        WriteXml
        WriteJsonError
        WriteXmlError
------------------------------------------------

      Name - logger
      Description - wrapper of uber/zap
      Functions - 
        Info
        Debug
        Error
------------------------------------------------

      Name - security
      Description - crypto/bcryct hashing
      Functions - 
        HashPassword
        VerifyPassword
        
