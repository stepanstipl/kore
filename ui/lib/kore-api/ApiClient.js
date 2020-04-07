/**
 * The auto-generated models depend on a class called ApiClient for some type conversions, this very-cut-down 
 * ApiClient provides just those static functions. Based on the auto-generated ApiClient.js with all unnecessary
 * functions removed.
 */
class ApiClient {
    /**
    * Parses an ISO-8601 string representation of a date value.
    * @param {String} str The date value as a string.
    * @returns {Date} The parsed date object.
    */
    static parseDate(str) {
        return new Date(str);
    }

    /**
    * Converts a value to the specified type.
    * @param {(String|Object)} data The data to convert, as a string or object.
    * @param {(String|Array.<String>|Object.<String, Object>|Function)} type The type to return. Pass a string for simple types
    * or the constructor function for a complex type. Pass an array containing the type name to return an array of that type. To
    * return an object, pass an object with one property whose name is the key type and whose value is the corresponding value type:
    * all properties on <code>data<code> will be converted to this type.
    * @returns An instance of the specified type or null or undefined if data is null or undefined.
    */
    static convertToType(data, type) {
        if (data === null || data === undefined)
            return data

        switch (type) {
            case 'Boolean':
                return Boolean(data);
            case 'Integer':
                return parseInt(data, 10);
            case 'Number':
                return parseFloat(data);
            case 'String':
                return String(data);
            case 'Date':
                return ApiClient.parseDate(String(data));
            case 'Blob':
                return data;
            default:
                if (type === Object) {
                    // generic object, return directly
                    return data;
                } else if (typeof type.constructFromObject === 'function') {
                    // for model type like User and enum class
                    return type.constructFromObject(data);
                } else if (Array.isArray(type)) {
                    // for array type like: ['String']
                    var itemType = type[0];

                    return data.map((item) => {
                        return ApiClient.convertToType(item, itemType);
                    });
                } else if (typeof type === 'object') {
                    // for plain object type like: {'String': 'Integer'}
                    var keyType, valueType;
                    for (var k in type) {
                        if (type.hasOwnProperty(k)) {
                            keyType = k;
                            valueType = type[k];
                            break;
                        }
                    }

                    var result = {};
                    for (var k in data) {
                        if (data.hasOwnProperty(k)) {
                            var key = ApiClient.convertToType(k, keyType);
                            var value = ApiClient.convertToType(data[k], valueType);
                            result[key] = value;
                        }
                    }

                    return result;
                } else {
                    // for unknown type, return the data directly
                    return data;
                }
        }
    }

    /**
    * Constructs a new map or array model from REST data.
    * @param data {Object|Array} The REST data.
    * @param obj {Object|Array} The target object or array.
    */
    static constructFromObject(data, obj, itemType) {
        if (Array.isArray(data)) {
            for (var i = 0; i < data.length; i++) {
                if (data.hasOwnProperty(i))
                    obj[i] = ApiClient.convertToType(data[i], itemType);
            }
        } else {
            for (var k in data) {
                if (data.hasOwnProperty(k))
                    obj[k] = ApiClient.convertToType(data[k], itemType);
            }
        }
    };
}

export default ApiClient;
