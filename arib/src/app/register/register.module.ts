import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { BrowserModule } from "@angular/platform-browser";
import { FormsModule } from "@angular/forms";

import { NgbModule } from "@ng-bootstrap/ng-bootstrap";
import { Ng2Webstorage } from "ng2-webstorage";

import { RegisterRoutingModule } from "./register-routing.module";

import { RegisterComponent } from "./register.component";

@NgModule({
	imports: [
		CommonModule,
		BrowserModule,
		FormsModule,

		NgbModule,
		Ng2Webstorage,

		RegisterRoutingModule
	],
	providers: [],
	declarations: [
		RegisterComponent
	]
})
export class RegisterModule { }
