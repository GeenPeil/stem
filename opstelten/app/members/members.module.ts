import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { CommonModule } from '../common/common.module';
import { ApiModule } from '../api/api.module';
import { MembersRoutingModule } from './members-routing.module';

import { MemberService } from './member.service';
import { MemberListComponent } from './member-list.component';
import { MemberComponent } from './member.component';

@NgModule({
	imports: [
		BrowserModule,
		FormsModule,
		NgbModule,
		CommonModule,
		ApiModule,
		MembersRoutingModule
	],
	providers: [
		MemberService
	],
	declarations: [
		MemberListComponent,
		MemberComponent
	]
})
export class MembersModule { }
