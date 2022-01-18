#Again, works fine directly under GetBash, but not thru batch file.
#Solution? Rename the .bat to .sh file, especially when is exuected under gitBash
curl -i -H "Content-Type: application/json" -X POST -d '{"Id":"13","Title":"Test Article 2 psuedo title","Date":"2022-01-17","Body":"Psuedo body","Tags": ["tag1", "tag2", "tag3"]}' http://localhost:8080/articles