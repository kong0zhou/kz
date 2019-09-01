import { Injectable } from '@angular/core';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import {ReplyProto,ReqProto} from "../msg-proto";
import { Observable ,of} from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class MainService {

  constructor(
    private http:HttpClient,
  ) { }

  startUp(d): Observable<ReplyProto> {
    let url='/ab-api/deploy'
    return this.http.put<ReplyProto>(url, d)
  }
}
