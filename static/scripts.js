window.UploadPDF= function UploadPDF(e) 
{
    //let user = { name:'john', age:34 };
    let xhr = new XMLHttpRequest();
    let formData = new FormData();
    let pdfSubmission = e.files[0];      
    
    //formData.append("user", JSON.stringify(user));   
    formData.append("pdf", pdfSubmission);
    
    xhr.onreadystatechange = state => { console.log(xhr.status); } // err handling
    xhr.open("POST", '/upload/pdf');    
    xhr.send(formData);
}
