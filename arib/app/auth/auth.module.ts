import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { CommonModule } from '../common/common.module';
import { AuthRoutingModule } from './auth-routing.module';
import { Auth } from './auth.service';

import { LoginComponent } from './login.component';

@NgModule({
	imports: [
		FormsModule,
		HttpModule,
		CommonModule,
		AuthRoutingModule
	],
	providers: [
		Auth
	],
	declarations: [
		LoginComponent
	]
})
export class AuthModule { }
