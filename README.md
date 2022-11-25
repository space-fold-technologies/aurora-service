## <B>SETUP</B>

- set up the a daemon to run the service or run it in a background service ``` 
./aurora-service --mode=run ``` this will also run migrations and set up all tables
- register a super administration user and initialize the default settings in your ```settings.yml``` run ```./aurora-service --mode=setup --email=admin@mail.domain --password=*************```
- to nuke the set up ```./aurora-service --mode=reset```
- still working out the issues in the code base


