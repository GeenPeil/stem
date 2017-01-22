import { NgModule } from '@angular/core';

import { ConfigService } from './config.service';

import { DateToStringPipe } from './date-to-string.pipe';

@NgModule({
    imports: [],
    providers: [ConfigService],
    declarations: [DateToStringPipe],
    exports: [DateToStringPipe]
})
export class CommonModule { }
