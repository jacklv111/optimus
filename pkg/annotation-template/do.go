/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package annotationtemplate

type AnnotationTemplateDo struct {
	ID            string `gorm:"primary_key;<-:create"`
	AssociationId string `gorm:"type:varchar(255)"`
}

func (AnnotationTemplateDo) TableName() string {
	return "annotation_templates"
}
