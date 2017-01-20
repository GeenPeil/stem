import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { AuthModule } from '../auth/auth.module';
import { ProfileRoutingModule } from './profile-routing.module';
import { AuthHttp } from '../auth/auth-http.service';

import { ProfileComponent } from './profile.component';

@NgModule({
    imports: [
        FormsModule,
        HttpModule,
        AuthModule,
        ProfileRoutingModule
    ],
    providers: [
        AuthHttp
    ],
    declarations: [
        ProfileComponent
    ]
})
export class ProfileModule { }
