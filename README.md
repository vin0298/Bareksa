Running the program might require to manually use the schema to create the
database (I apologise in advance for any inconvenience).  

Running  
 ```go build main.go```  
then  
```./bareksa```  
should be sufficient to allow the program to run.

## CONCERNS/THOUGHTS
I noticed when reviewing that there are some cases (especially in the repository) that I might have been violating the DRY principle. 

Moreover, I think that the NewsArticleRepository is overburdened and there should be a separate handler for Tags. The reason I didn't do so it's because of my reasoning that Tags are a part of a NewsArticle - it doesn't really makes sense to split it apart. 

Lack of proper testing 

# PRIORITY:
COMMENTS + README  
INTEGRATION TESTING

## BUGS:
Updating an article with the same array of tags causes a SQL syntax error at deleteUnusedTagsJoinTable  
Migration is not working yet.
