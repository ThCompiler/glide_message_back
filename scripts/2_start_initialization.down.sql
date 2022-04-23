WITH json (doc) AS (
    values
        ('[
          {
            "language": "Australia",
            "picture": "."
          },
          {
            "language": "Armenia",
            "picture": "."
          },
          {
            "language": "Belarus",
            "picture": "."
          },
          {
            "language": "Brazil",
            "picture": "."
          },
          {
            "language": "Canada",
            "picture": "."
          },
          {
            "language": "France",
            "picture": "."
          },
          {
            "language": "Germany",
            "picture": "."
          },
          {
            "language": "India",
            "picture": "."
          },
          {
            "language": "Japan",
            "picture": "."
          },
          {
            "language": "Republic_of_China",
            "picture": "."
          },
          {
            "language": "Russia",
            "picture": "."
          },
          {
            "language": "Spain",
            "picture": "."
          },
          {
            "language": "Switzerland",
            "picture": "."
          },
          {
            "language": "Turkey",
            "picture": "."
          },
          {
            "language": "UK",
            "picture": "."
          },
          {
            "language": "Ukraine",
            "picture": "."
          },
          {
            "language": "USA",
            "picture": "."
          },
          {
            "language": "Vietnam",
            "picture": "."
          }
        ]'::json)
)
DELETE FROM countries WHERE country_name IN (
SELECT p.country_name
from json l cross join lateral json_populate_recordset(null::countries, doc) as p);

WITH json (doc) AS (
    values
        ('[
          {
            "country_name": "Chinese",
            "picture": "."
          },
          {
            "country_name": "English",
            "picture": "."
          },
          {
            "country_name": "Spanish",
            "picture": "."
          },
          {
            "country_name": "Arabian",
            "picture": "."
          },
          {
            "country_name": "Russian",
            "picture": "."
          },
          {
            "country_name": "Portuguese",
            "picture": "."
          },
          {
            "country_name": "French",
            "picture": "."
          },
          {
            "country_name": "German",
            "picture": "."
          },
          {
            "country_name": "Hindi",
            "picture": "."
          },
          {
            "country_name": "Bengali",
            "picture": "."
          },
          {
            "country_name": "Japanese",
            "picture": "."
          },
          {
            "country_name": "Italian",
            "picture": "."
          },
          {
            "country_name": "Belarusian",
            "picture": "."
          },
          {
            "country_name": "Ukrainian",
            "picture": "."
          },
          {
            "country_name": "Vietnamese",
            "picture": "."
          },
          {
            "country_name": "Armenian",
            "picture": "."
          },
          {
            "country_name": "Turkish",
            "picture": "."
          }
        ]'::json)
)
DELETE FROM languages WHERE language IN (
    SELECT p.language
    from json l cross join lateral json_populate_recordset(null::languages, doc) as p);



