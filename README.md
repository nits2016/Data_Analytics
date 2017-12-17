# Data_Analytics

The aim to is to visualise the given data which contains the age-groups(in range) with respect to the drugs consumption. The data sheet is the authentic data of the Indiana state of USA.
Required software needed -> PostgreSQL, metabase for visualisation, GO language, Superset for detailed Visualization.

# Steps
1. The dataset used here can be found here :
PostgreSQL database is used for storing the data(table). We can use any database which are supported by the metabase and superset. The import of the data from the spreadsheet or sql(file) to postgreSQL has been done by the json file and the GO adapter attached in the project. 
Firstly check the dataset i.e .csv file or check the schema of the database(.sql file is provided to directly load the data).Create a database and inside it similar kind of table in postgres( use pgadmin4). The table has one column i.e AGE - it comes in the datatype INT4RANGE. You can directly import the csv file into the table or directly run the go adaptor file( ).

