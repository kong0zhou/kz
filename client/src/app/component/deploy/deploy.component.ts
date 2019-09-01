import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroupDirective, NgForm, Validators } from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import { MainService } from '../../services/main.service';
import { ReplyProto, ReqProto } from "../../msg-proto";

@Component({
  selector: 'app-deploy',
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.css']
})
export class DeployComponent implements OnInit {
  // =======验证器========
  bfFormControl = new FormControl('', [
    Validators.required,
    Validators.pattern(`[0-9]+`),
  ]);
  timeFormControl = new FormControl('', [
    Validators.required,
    Validators.pattern(`[0-9]+`),
  ]);
  urlFormControl = new FormControl('', [
    Validators.required,
  ])
  qDataFormControl = new FormControl('', [
    Validators.required,
  ])
  pathFormControl = new FormControl('', [
    Validators.required,
  ])

  // =======input变量====
  bf = ''// 并法数
  method = 'POST'; // 请求方法
  qData = ""// 请求携带数据
  path = "" // 导出结果路径
  time = ""
  timeUnit = "分钟"
  url = ""


  constructor(
    public main: MainService,
  ) { }

  ngOnInit() {
  }

  startUp() {
    let req: ReqProto = {
      data: {
        concurrent: Number(this.bf),
        method: this.method,
        reqData: this.qData,
        filePath: this.path,
        time: this.convertTime(),
        url: this.url
      }
    }
    console.log(JSON.stringify(req))
    this.main.startUp(req).subscribe(data => {
      console.log(JSON.stringify(data))
    },
      (error) => {
        console.log(JSON.stringify(error))
      }
    )
  }

  convertTime(): number {
    switch (this.timeUnit) {
      case "小时":
        return Number(this.time) * 3600;
      case "分钟":
        return Number(this.time) * 60;
      case "秒":
        return Number(this.time)
      default:
        console.log("convertTime has something bad happened")
    }
  }

  onKeyInfo(value, type) {
    if (typeof value == 'undefined' || value == null || value == '') {
      console.log("error:value is null")
    }
    if (typeof type == 'undefined' || type == null || value == '') {
      console.log("error:type is null")
    }
    switch (type) {
      case 'url':
        this.url = value;
        break;
      case 'qData':
        this.qData = value;
        break;
      case 'path':
        this.path = value;
        break;
      case 'bf':
        this.bf = value;
        break;
      case 'time':
        this.time = value;
        break;
      default:
        console.log("error:type is no found")
    }
  }
}
