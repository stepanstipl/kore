/**
 * Kore API
 * Kore API provides the frontend API (kore.appvia.io)
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
 * The V1SecretReference model module.
 * @module model/V1SecretReference
 * @version 0.0.1
 */
class V1SecretReference {
    /**
     * Constructs a new <code>V1SecretReference</code>.
     * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
     * @alias module:model/V1SecretReference
     */
    constructor() { 
        
        V1SecretReference.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1SecretReference</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1SecretReference} obj Optional instance to populate.
     * @return {module:model/V1SecretReference} The populated <code>V1SecretReference</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1SecretReference();

            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('namespace')) {
                obj['namespace'] = ApiClient.convertToType(data['namespace'], 'String');
            }
        }
        return obj;
    }

/**
     * Returns Name is unique within a namespace to reference a secret resource.
     * @return {String}
     */
    getName() {
        return this.name;
    }

    /**
     * Sets Name is unique within a namespace to reference a secret resource.
     * @param {String} name Name is unique within a namespace to reference a secret resource.
     */
    setName(name) {
        this['name'] = name;
    }
/**
     * Returns Namespace defines the space within which the secret name must be unique.
     * @return {String}
     */
    getNamespace() {
        return this.namespace;
    }

    /**
     * Sets Namespace defines the space within which the secret name must be unique.
     * @param {String} namespace Namespace defines the space within which the secret name must be unique.
     */
    setNamespace(namespace) {
        this['namespace'] = namespace;
    }

}

/**
 * Name is unique within a namespace to reference a secret resource.
 * @member {String} name
 */
V1SecretReference.prototype['name'] = undefined;

/**
 * Namespace defines the space within which the secret name must be unique.
 * @member {String} namespace
 */
V1SecretReference.prototype['namespace'] = undefined;






export default V1SecretReference;

