import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { CommonModule } from '../common/common.module';
import { AuthHttp } from './auth-http.service';


@NgModule({
    imports: [
        HttpModule,
        CommonModule
    ],
    providers: [
        AuthHttp
    ],
    declarations: []
})
export class AuthModule { }
