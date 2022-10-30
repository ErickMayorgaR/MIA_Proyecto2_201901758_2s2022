import { Component, Input } from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http'
import { Observable } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  constructor(private http:HttpClient) {}
  

  textBox1:string = "";
  url:string = "http://localhost:9090/";
  response:string = "";

  

  getInformation(){
    this.makeRequest().subscribe(apiResponse =>{
      this.response = apiResponse;  
    })
  }


  makeRequest():Observable<any>{
    
    let address = this.url + "get"
    return this.http.get(address,{responseType: 'text'})
  }

  }
  



