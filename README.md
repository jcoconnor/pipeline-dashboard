# Pipeline Dashboard

The pipeline dashboard is a tool for measuring the length of time that Jenkins builds with deep dependencies take to run, and then number of errors in those jobs.

# Running

Copy conf/config.example.toml to conf/config.toml and add your Product and Jenkins Job information there.  A "product" is a codebase which may have multiple branches being monitored, with multiple Jenkins Pipelines

To get data to show in the dashboard, run ```main.go```.  To run the API, run ```cmd/web/main.go```  In order to run the react frontend, cd into frontend and run ```yarn start```.
