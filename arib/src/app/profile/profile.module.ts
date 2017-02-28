import { NgModule } from "@angular/core";

import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";
import { HttpModule } from "@angular/http";

import { ApiModule } from "../api/api.module";
import { ProfileRoutingModule } from "./profile-routing.module";

import { ProfileComponent } from "./profile.component";

@NgModule({
	imports: [
		CommonModule,
		FormsModule,
		HttpModule,

		ApiModule,
		ProfileRoutingModule
	],
	providers: [],
	declarations: [
		ProfileComponent
	]
})
export class ProfileModule { }
