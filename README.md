# Data_Analytics

The aim is to visualise the given data which contains the age-groups(in range) with respect to the drugs consumption. The data sheet is the authentic data of the Indiana state of USA.
Required software needed -> PostgreSQL, metabase for visualisation, GO language, Superset for detailed Visualization.

# Steps
1. The dataset used here can be found here : dataset.sql 
#

PostgreSQL database is used for storing the data(table). We can use any database which are supported by the metabase and superset. The import of the data from the spreadsheet or sql(file convertcsv.sql) to postgreSQL has been done by the json file and the GO adapter attached in the project. 
Firstly check the dataset i.e .csv file(drug-use-by-age.csv) or check the schema of the database(.sql file is provided to directly load the data).Create a database and inside it similar kind of table in postgres( use pgadmin4). The table has one column i.e AGE - it comes in the datatype INT4RANGE. You can directly import the csv file into the table or directly run the go adaptor file( ). The main advantage of using the GO adaptor is that it can import multiple tables in just one program. Moreover GO language doesn't consist problems like memory exceeded, heap memory, etc.
The adaptor works as follows:- We have made a JSON file in which we map the columns with the csv file. The GO script then import all the file data into the postgres Database.
Always run the "TableLoadThroughMap.exe" in cmd to easily observe the errors.Dont't forget to remove all comments from the "config_TableMap_Bank1.json".
Various errors and solutions while running the Adaptor file
1. {{ {{}} }} empty set. This means that there must be some error in the file path written in "config_TableMap_Bank1.json". Try using double \\ if your file is in C drive or single / if your file is in other drive.
2. Try putting the correct table name , databse name, etc.
3. Check that your table in postgres contains same amount of column as in the json file.

# NOTE
Before running "TableLoadThroughMap.exe" remove all comments from the "config_TableMap_Bank1.json".


2. After completing the database part we come to the visualisation of data. Open your metabase and then open the admin panel. Add a databse with all the details and correct password. Now the newly added database will be visible. After adding the data its very important to again go into the showing database and click sync database schema now.

3. The various visualisations and results are provided in Results.pdf. Do use the real images i.e 1.png,2.png,etc to clearly observe the visualisation.

