Dengan go, buatlah project API posting board sederhana, hanya API, tidak perlu frontend,
API Response nya dalam format JSON

- menggunakan docker untuk menjalankan program API nya,  
- Authentikasi dengan JWT  
- Thread system
  - Create thread, containing title and text
  - Delete thread 
  - List own thread 
  - List public thread
    - Sort by most/least commented thread
    - Sort by most/least liked thread
    - Find post by poster user id
    - Find post that a user has comment on
- User bisa meninggalkan komentar di post
  - post comment, containing text only
  - view post comment
- User can like a thread
  - Melihat list user yang like post tersebut

Note:  
Implement error handling for unexpected exceptions and expected exception  
Validation only implement in one of API of your choice  
Using gorm to auto migrate and as your database ORM if possible

