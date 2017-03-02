import { NgModule } from "@angular/core";
import { HttpModule } from "@angular/http";

import { Common } from "../common/common.module";
import { AuthModule } from "../auth/auth.module";

import { Api } from "./api.service";

@NgModule({
	imports: [
		HttpModule,

		Common,
		AuthModule
	],
	providers: [Api],
	declarations: []
})
export class ApiModule { }
