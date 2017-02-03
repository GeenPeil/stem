import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

import { AuthModule } from '../auth/auth.module';
import { ProfileRoutingModule } from './profile-routing.module';

import { ProfileComponent } from './profile.component';

@NgModule({
    imports: [
        FormsModule,
        CommonModule,
        HttpModule,
        AuthModule,
        ProfileRoutingModule
    ],
    providers: [],
    declarations: [
        ProfileComponent
    ]
})
export class ProfileModule { }
