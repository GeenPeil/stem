import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

import { Auth } from '../auth/auth.service';
import { AuthModule } from '../auth/auth.module';
import { ProfileRoutingModule } from './profile-routing.module';
import { AuthHttp } from '../auth/auth-http.service';

import { ProfileComponent } from './profile.component';

@NgModule({
    imports: [
        FormsModule,
        CommonModule,
        HttpModule,
        AuthModule,
        ProfileRoutingModule
    ],
    providers: [
        Auth,
        AuthHttp
    ],
    declarations: [
        ProfileComponent
    ]
})
export class ProfileModule { }
