import { NgModule } from '@angular/core';

import { CommonModule } from '@angular/common';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { Common } from '../common/common.module';
import { AuthRoutingModule } from './auth-routing.module';
import { Auth } from './auth.service';

import { LoginComponent } from './login.component';
import { LogoutComponent } from './logout.component';

@NgModule({
	imports: [
		CommonModule,
		FormsModule,
		HttpModule,

		Common,

		AuthRoutingModule
	],
	providers: [
		Auth
	],
	declarations: [
		LoginComponent,
		LogoutComponent
	]
})
export class AuthModule { }
