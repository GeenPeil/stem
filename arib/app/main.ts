import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app.module';

var appPromise = platformBrowserDynamic().bootstrapModule(AppModule);

