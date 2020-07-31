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
 * The V1beta1CostAsset model module.
 * @module model/V1beta1CostAsset
 * @version 0.0.1
 */
class V1beta1CostAsset {
    /**
     * Constructs a new <code>V1beta1CostAsset</code>.
     * @alias module:model/V1beta1CostAsset
     */
    constructor() { 
        
        V1beta1CostAsset.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>V1beta1CostAsset</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/V1beta1CostAsset} obj Optional instance to populate.
     * @return {module:model/V1beta1CostAsset} The populated <code>V1beta1CostAsset</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new V1beta1CostAsset();

            if (data.hasOwnProperty('assetIdentifier')) {
                obj['assetIdentifier'] = ApiClient.convertToType(data['assetIdentifier'], 'String');
            }
            if (data.hasOwnProperty('koreIdentifier')) {
                obj['koreIdentifier'] = ApiClient.convertToType(data['koreIdentifier'], 'String');
            }
            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('provider')) {
                obj['provider'] = ApiClient.convertToType(data['provider'], 'String');
            }
            if (data.hasOwnProperty('tags')) {
                obj['tags'] = ApiClient.convertToType(data['tags'], {'String': 'String'});
            }
            if (data.hasOwnProperty('teamIdentifier')) {
                obj['teamIdentifier'] = ApiClient.convertToType(data['teamIdentifier'], 'String');
            }
        }
        return obj;
    }

/**
     * @return {String}
     */
    getAssetIdentifier() {
        return this.assetIdentifier;
    }

    /**
     * @param {String} assetIdentifier
     */
    setAssetIdentifier(assetIdentifier) {
        this['assetIdentifier'] = assetIdentifier;
    }
/**
     * @return {String}
     */
    getKoreIdentifier() {
        return this.koreIdentifier;
    }

    /**
     * @param {String} koreIdentifier
     */
    setKoreIdentifier(koreIdentifier) {
        this['koreIdentifier'] = koreIdentifier;
    }
/**
     * @return {String}
     */
    getName() {
        return this.name;
    }

    /**
     * @param {String} name
     */
    setName(name) {
        this['name'] = name;
    }
/**
     * @return {String}
     */
    getProvider() {
        return this.provider;
    }

    /**
     * @param {String} provider
     */
    setProvider(provider) {
        this['provider'] = provider;
    }
/**
     * @return {Object.<String, String>}
     */
    getTags() {
        return this.tags;
    }

    /**
     * @param {Object.<String, String>} tags
     */
    setTags(tags) {
        this['tags'] = tags;
    }
/**
     * @return {String}
     */
    getTeamIdentifier() {
        return this.teamIdentifier;
    }

    /**
     * @param {String} teamIdentifier
     */
    setTeamIdentifier(teamIdentifier) {
        this['teamIdentifier'] = teamIdentifier;
    }

}

/**
 * @member {String} assetIdentifier
 */
V1beta1CostAsset.prototype['assetIdentifier'] = undefined;

/**
 * @member {String} koreIdentifier
 */
V1beta1CostAsset.prototype['koreIdentifier'] = undefined;

/**
 * @member {String} name
 */
V1beta1CostAsset.prototype['name'] = undefined;

/**
 * @member {String} provider
 */
V1beta1CostAsset.prototype['provider'] = undefined;

/**
 * @member {Object.<String, String>} tags
 */
V1beta1CostAsset.prototype['tags'] = undefined;

/**
 * @member {String} teamIdentifier
 */
V1beta1CostAsset.prototype['teamIdentifier'] = undefined;






export default V1beta1CostAsset;

