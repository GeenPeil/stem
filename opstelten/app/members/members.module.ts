import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { CommonModule } from '../common/common.module';
import { AuthModule } from '../auth/auth.module';
import { MembersRoutingModule } from './members-routing.module';

import { MemberService } from './member.service';

import { MemberComponent } from './member.component';

@NgModule({
    imports: [
        BrowserModule,
        FormsModule,
        NgbModule,
        CommonModule,
        AuthModule,
        MembersRoutingModule
    ],
    providers: [
        MemberService
    ],
    declarations: [MemberComponent]
})
export class MembersModule { }
