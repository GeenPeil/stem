import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { CommonModule } from '../common/common.module';
import { AuthModule } from '../auth/auth.module';

import { Api } from './api.service';


@NgModule({
	imports: [
		HttpModule,
		CommonModule,
		AuthModule
	],
	providers: [Api],
	declarations: []
})
export class ApiModule { }
