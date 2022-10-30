import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http'




export class InformationAPIService {

  url:string = "http://localhost:4200/";

  constructor(private http:HttpClient) { }

  getInformation(){
    return this.http.get("get");
  }


}
