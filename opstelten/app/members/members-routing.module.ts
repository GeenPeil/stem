import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MemberComponent } from './member.component';

const membersRoutes: Routes = [
    { path: 'member/:id', component: MemberComponent }
];

@NgModule({
    imports: [
        RouterModule.forChild(membersRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class MembersRoutingModule { }
