/**
 * Appvia Kore API
 * Kore API provides the frontend API for the Appvia Kore (kore.appvia.io)
 *
 * The version of the OpenAPI document: 0.0.1
 * Contact: info@appvia.io
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The KoreServicePlanDetails model module.
 * @module model/KoreServicePlanDetails
 * @version 0.0.1
 */
class KoreServicePlanDetails {
    /**
     * Constructs a new <code>KoreServicePlanDetails</code>.
     * @alias module:model/KoreServicePlanDetails
     * @param editableParams {Array.<String>} 
     * @param kind {String} 
     * @param summary {String} 
     */
    constructor(editableParams, kind, summary) { 
        
        KoreServicePlanDetails.initialize(this, editableParams, kind, summary);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, editableParams, kind, summary) { 
        obj['editableParams'] = editableParams;
        obj['kind'] = kind;
        obj['summary'] = summary;
    }

    /**
     * Constructs a <code>KoreServicePlanDetails</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/KoreServicePlanDetails} obj Optional instance to populate.
     * @return {module:model/KoreServicePlanDetails} The populated <code>KoreServicePlanDetails</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new KoreServicePlanDetails();

            if (data.hasOwnProperty('configuration')) {
                obj['configuration'] = ApiClient.convertToType(data['configuration'], 'String');
            }
            if (data.hasOwnProperty('credentialSchema')) {
                obj['credentialSchema'] = ApiClient.convertToType(data['credentialSchema'], 'String');
            }
            if (data.hasOwnProperty('description')) {
                obj['description'] = ApiClient.convertToType(data['description'], 'String');
            }
            if (data.hasOwnProperty('displayName')) {
                obj['displayName'] = ApiClient.convertToType(data['displayName'], 'String');
            }
            if (data.hasOwnProperty('editableParams')) {
                obj['editableParams'] = ApiClient.convertToType(data['editableParams'], ['String']);
            }
            if (data.hasOwnProperty('kind')) {
                obj['kind'] = ApiClient.convertToType(data['kind'], 'String');
            }
            if (data.hasOwnProperty('labels')) {
                obj['labels'] = ApiClient.convertToType(data['labels'], {'String': 'String'});
            }
            if (data.hasOwnProperty('providerData')) {
                obj['providerData'] = ApiClient.convertToType(data['providerData'], 'String');
            }
            if (data.hasOwnProperty('schema')) {
                obj['schema'] = ApiClient.convertToType(data['schema'], 'String');
            }
            if (data.hasOwnProperty('serviceAccessDisabled')) {
                obj['serviceAccessDisabled'] = ApiClient.convertToType(data['serviceAccessDisabled'], 'Boolean');
            }
            if (data.hasOwnProperty('summary')) {
                obj['summary'] = ApiClient.convertToType(data['summary'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getConfiguration() {
        return this.configuration;
    }

    /**
     * @param {String} configuration
     */
    setConfiguration(configuration) {
        this['configuration'] = configuration;
    }
/**
     * @return {String}
     */
    getCredentialSchema() {
        return this.credentialSchema;
    }

    /**
     * @param {String} credentialSchema
     */
    setCredentialSchema(credentialSchema) {
        this['credentialSchema'] = credentialSchema;
    }
/**
     * @return {String}
     */
    getDescription() {
        return this.description;
    }

    /**
     * @param {String} description
     */
    setDescription(description) {
        this['description'] = description;
    }
/**
     * @return {String}
     */
    getDisplayName() {
        return this.displayName;
    }

    /**
     * @param {String} displayName
     */
    setDisplayName(displayName) {
        this['displayName'] = displayName;
    }
/**
     * @return {Array.<String>}
     */
    getEditableParams() {
        return this.editableParams;
    }

    /**
     * @param {Array.<String>} editableParams
     */
    setEditableParams(editableParams) {
        this['editableParams'] = editableParams;
    }
/**
     * @return {String}
     */
    getKind() {
        return this.kind;
    }

    /**
     * @param {String} kind
     */
    setKind(kind) {
        this['kind'] = kind;
    }
/**
     * @return {Object.<String, String>}
     */
    getLabels() {
        return this.labels;
    }

    /**
     * @param {Object.<String, String>} labels
     */
    setLabels(labels) {
        this['labels'] = labels;
    }
/**
     * @return {String}
     */
    getProviderData() {
        return this.providerData;
    }

    /**
     * @param {String} providerData
     */
    setProviderData(providerData) {
        this['providerData'] = providerData;
    }
/**
     * @return {String}
     */
    getSchema() {
        return this.schema;
    }

    /**
     * @param {String} schema
     */
    setSchema(schema) {
        this['schema'] = schema;
    }
/**
     * @return {Boolean}
     */
    getServiceAccessDisabled() {
        return this.serviceAccessDisabled;
    }

    /**
     * @param {Boolean} serviceAccessDisabled
     */
    setServiceAccessDisabled(serviceAccessDisabled) {
        this['serviceAccessDisabled'] = serviceAccessDisabled;
    }
/**
     * @return {String}
     */
    getSummary() {
        return this.summary;
    }

    /**
     * @param {String} summary
     */
    setSummary(summary) {
        this['summary'] = summary;
    }

}

/**
 * @member {String} configuration
 */
KoreServicePlanDetails.prototype['configuration'] = undefined;

/**
 * @member {String} credentialSchema
 */
KoreServicePlanDetails.prototype['credentialSchema'] = undefined;

/**
 * @member {String} description
 */
KoreServicePlanDetails.prototype['description'] = undefined;

/**
 * @member {String} displayName
 */
KoreServicePlanDetails.prototype['displayName'] = undefined;

/**
 * @member {Array.<String>} editableParams
 */
KoreServicePlanDetails.prototype['editableParams'] = undefined;

/**
 * @member {String} kind
 */
KoreServicePlanDetails.prototype['kind'] = undefined;

/**
 * @member {Object.<String, String>} labels
 */
KoreServicePlanDetails.prototype['labels'] = undefined;

/**
 * @member {String} providerData
 */
KoreServicePlanDetails.prototype['providerData'] = undefined;

/**
 * @member {String} schema
 */
KoreServicePlanDetails.prototype['schema'] = undefined;

/**
 * @member {Boolean} serviceAccessDisabled
 */
KoreServicePlanDetails.prototype['serviceAccessDisabled'] = undefined;

/**
 * @member {String} summary
 */
KoreServicePlanDetails.prototype['summary'] = undefined;






export default KoreServicePlanDetails;
